<script lang="ts">
	import '../app.postcss';
	import { popup } from '@skeletonlabs/skeleton';
	import type { PopupSettings, ToastSettings } from '@skeletonlabs/skeleton';
	import { Api, type HttpResponse, type MainCreateLobbyRequest, type MainJoinLobbyRequest } from '../dtmd_api';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
    import { storePopup } from '@skeletonlabs/skeleton';
	import { LobbyID, MemberID,ErrorMessageStart } from '../stores.js';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { initializeStores, Toast } from '@skeletonlabs/skeleton';
	import { getToastStore } from '@skeletonlabs/skeleton';
	//inbuilt function from skeleton necessary for initialization
	initializeStores();
	const toastStore = getToastStore();
    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

    //Set base url for api to dev or prod 
    const isServerSide = typeof window === 'undefined';
    const apiBase = new Api({
        baseUrl: isServerSide ? process.env.VITE_API_BASE_URL : import.meta.env.VITE_API_BASE_URL
    });
	const api = apiBase.api;

	const popupCreateRoom: PopupSettings = {
		// Represents the type of event that opens/closed the popup
		event: 'click',
		// Matches the data-popup value on your popup element
		target: 'popupCreateRoom',
		// Defines which side of your trigger the popup will appear
		placement: 'bottom',
		//Defines which element closes poupup
		closeQuery: ''
	};
	const popupJoinRoom: PopupSettings = {
		// Represents the type of event that opens/closed the popup
		event: 'click',
		// Matches the data-popup value on your popup element
		target: 'popupJoinRoom',
		// Defines which side of your trigger the popup will appear
		placement: 'bottom',
		//Defines which element closes poupup
		closeQuery: ''
	};

	let roomname = '';
	let usernameCreate = '';
	let roomPin = '';
	let usernameJoin = '';
	let errorMessageCreate = '';
	let errorMessageJoin = '';
	let disableCreateButton = false;
	let disableJoinButton = false;
	onMount(() => {
		//Check if we return from different page with error
		const unsubscribeErrorMessageStart = ErrorMessageStart.subscribe(value => {
			if(value != null && value !==''){
				const t: ToastSettings = {
					message: "Lobby no longer exists.",
					timeout: 5000,
					background: 'variant-filled-error'
				};
				toastStore.trigger(t);
				ErrorMessageStart.set("");
			}
		});
		unsubscribeErrorMessageStart();
	});
	async function createRoom(tmpRoomname:string,tmpUsernameCreate:string) {
		//Check input
		if(tmpRoomname === ''){
			errorMessageCreate = "Please enter a room name.";
			setTimeout(() => {
                errorMessageCreate = '';
            }, 3000);
			return;
		}
		if(tmpUsernameCreate === ''){
			errorMessageCreate = "Please enter a nickname.";
			setTimeout(() => {
                errorMessageCreate = '';
            }, 3000);
			return;
		}
		disableCreateButton=true;
		try {
			var pin = await createLobbyAPICall(tmpRoomname);
			console.log("Lobby Created",pin);
   		} catch (error) {
			console.error('Error creating lobby:', error);
			const t: ToastSettings = {
				message: "Couldn't create room.",
				timeout: 3000,
				hideDismiss: true,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
			disableCreateButton=false;
			return;
		}
		try {
			var userID = await joinLobbyAPICall(pin,tmpUsernameCreate);
			console.log("Joined created lobby",pin);
   		} catch (error) {
			console.error('Error joining created lobby:', error);
			const t: ToastSettings = {
				message: "Couldn't join created room: "+ pin+".",
				timeout: 3000,
				hideDismiss: true,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
			disableCreateButton=false;
			return;
		}
		MemberID.set(userID);
		LobbyID.set(pin);
		console.log("successfull creation",userID);
		goto('/lobby').then(() => {
			disableCreateButton=false;
		});
    }
	async function joinRoom(tmpRoomPin:string, tmpUsernameJoin:string){
		//Check input
		if(tmpRoomPin === ''){
			errorMessageJoin = "Please enter a room-pin.";
			setTimeout(() => {
				errorMessageJoin = '';
			}, 3000);
			return;
		}
		const validPin = /^[A-Z0-9]{6}$/;
		if(!validPin.test(tmpRoomPin)){
			errorMessageJoin = "Please enter a valid room pin.";
			setTimeout(() => {
				errorMessageJoin = '';
			}, 3000);
			return;
		}
		if(tmpUsernameJoin === ''){
			errorMessageJoin = "Please enter a nickname.";
			setTimeout(() => {
                errorMessageJoin = '';
            }, 3000);
			return;
		}
		disableJoinButton=true;
		try {
			var userID = await joinLobbyAPICall(tmpRoomPin,tmpUsernameJoin);
			console.log("Joined lobby",tmpRoomPin);
   		} catch (error) {
			var errorCode = error as number;
			console.error('Error joining lobby:', error);
			//handle case if room doesn't exist
			if(errorCode === 404 ){
				const t: ToastSettings = {
					message: "Room with pin: "+ tmpRoomPin+" doesn't exist.",
					timeout: 3000,
					hideDismiss: true,
					background: 'variant-filled-error'
				};
				toastStore.trigger(t);
				disableJoinButton=false;
				return;
			}
			const t: ToastSettings = {
				message: "Couldn't join room with pin: "+ tmpRoomPin+".",
				timeout: 3000,
				hideDismiss: true,
				background: 'variant-filled-error'
			};
			toastStore.trigger(t);
			disableJoinButton=false;
			return;
		}
		MemberID.set(userID);
		LobbyID.set(tmpRoomPin);
		console.log("successful lobby join",userID);
		goto('/lobby').then(() => {
			disableJoinButton=false;
		});
	}

	async function createLobbyAPICall(lobbyName:string): Promise<string>{
		const lobbyCreateRequest: MainCreateLobbyRequest = {
    		name: lobbyName
		};
		const res = await api.lobbiesCreate(lobbyCreateRequest);
		if (res.ok) {
			const lobbyPin =  await res.text();
			return lobbyPin;
		} else {
			throw new Error("Failed to create lobby.");
		}
	}
	async function joinLobbyAPICall(pin:string,username:string): Promise<string>{
		const lobbyJoinREquest: MainJoinLobbyRequest = {
    		nickname: username
		};
		//api throws error if http response not ok
		try {
			const res = await api.lobbiesMembersCreate(pin,lobbyJoinREquest)
			if (!res.ok) {
				throw new Error(res.statusText);
			}
			var userID =  await res.text();
			return userID;
		} catch (error) {
			const httpResponse = error as HttpResponse<any, any>;
			throw httpResponse.status;
		}
	}
</script>
<Toast />
<div class="flex justify-center items-center mt-5"><img src="logo.png" alt="logo" style="height: 100px; width: auto;"></div>

<h1 class="text-center text-8xl mt-3" style="font-weight: 900;">DTMD</h1>
<h3 class="text-center text-2xl mt-1" style="color: #D4163C;">Don't Touch My Dice</h3>
<div class="flex flex-col items-center justify-center text-center">
	<button
		class="btn variant-filled"
		style="margin-top: 20vh; font-size: 20px;"
		use:popup={popupCreateRoom}>Create Room</button
	>
	<button
	class="btn variant-filled"
	style="margin-top: 30px; font-size: 20px;"
	use:popup={popupJoinRoom}>Join Room</button>

</div>

<div class="fixed bottom-5 w-full text-center">Welcome to DTMD, an online platform for virtual dice games!<br>Create a room, invite a few of your friends and get rolling!</div>


<div class="card p-4 w-72 shadow-xl" data-popup="popupCreateRoom">
	<div class="flex flex-col">
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={roomname} placeholder="Room Name" />
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={usernameCreate} placeholder="Nickname" />
		{#if errorMessageCreate}
		<p class="text-red-500 mt-2">{errorMessageCreate}</p>
		{/if}
		<button class="btn variant-filled" style="margin-top: 2vh;" on:click={() => createRoom(roomname,usernameCreate)} disabled = {disableCreateButton} >Create</button>
	</div>
</div>

<div class="card p-4 w-72 shadow-xl" data-popup="popupJoinRoom">
	<div class="flex flex-col">
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={roomPin} placeholder="Room-Pin" />
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={usernameJoin} placeholder="Nickname" />
		{#if errorMessageJoin}
		<p class="text-red-500 mt-2">{errorMessageJoin}</p>
		{/if}
		<button class="btn variant-filled" style="margin-top: 2vh;" on:click={()=>joinRoom(roomPin,usernameJoin)} disabled = {disableJoinButton} >Join</button>
	</div>
</div>
