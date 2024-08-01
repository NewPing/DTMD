<script lang="ts">
	import '../app.postcss';
	import { popup } from '@skeletonlabs/skeleton';
	import type { PopupSettings } from '@skeletonlabs/skeleton';
	import { Api, type MainCreateLobbyRequest, type MainJoinLobbyRequest } from '../dtmd_api';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
    import { storePopup } from '@skeletonlabs/skeleton';
	import { LobbyID, MemberID } from '../stores.js';
	import { goto } from '$app/navigation';
    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

	const api = new Api({
		baseUrl: "https://diceapi.odysseyinvision.com",
	});

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

	async function createRoom() {
		//Check input
		if(roomname === ''){
			errorMessageCreate = "Please enter a room name.";
			setTimeout(() => {
                errorMessageCreate = '';
            }, 3000);
			return;
		}
		if(usernameCreate === ''){
			errorMessageCreate = "Please enter a nickname.";
			setTimeout(() => {
                errorMessageCreate = '';
            }, 3000);
			return;
		}
		disableCreateButton=true;
		var pin = await createLobbyAPICall(roomname);
		var userID = await joinLobbyAPICall(pin,usernameCreate);
		MemberID.set(userID);
		LobbyID.set(pin);
		console.log("successfull creation",userID);
		goto('/lobby').then(() => {
			disableCreateButton=false;
		});
    }
	async function joinRoom(){
		//Check input
		if(roomPin === ''){
			errorMessageJoin = "Please enter a room-pin.";
			setTimeout(() => {
				errorMessageJoin = '';
			}, 3000);
			return;
		}
		const validPin = /^[A-Z0-9]{6}$/;
		if(!validPin.test(roomPin)){
			errorMessageJoin = "Please enter a valid room pin.";
			setTimeout(() => {
				errorMessageJoin = '';
			}, 3000);
			return;
		}
		if(usernameJoin === ''){
			errorMessageJoin = "Please enter a nickname.";
			setTimeout(() => {
                errorMessageJoin = '';
            }, 3000);
			return;
		}
		disableJoinButton=true;
		var userID = await joinLobbyAPICall(roomPin,usernameJoin);
		MemberID.set(userID);
		LobbyID.set(roomPin);
		console.log("successful lobby join",userID);
		goto('/lobby').then(() => {
			disableJoinButton=false;
		});
	}

	async function createLobbyAPICall(lobbyName:string): Promise<string>{
		const lobbyCreateRequest: MainCreateLobbyRequest = {
    		name: lobbyName
		};
		const res = await api.lobbies.lobbiesCreate(lobbyCreateRequest);
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
		const res = await api.lobbies.membersCreate(pin,lobbyJoinREquest)
		if (res.ok) {
			var userID =  await res.text();
			return userID;
		} else {
			throw new Error("Failed to join lobby.");
		}
	}
</script>

<h1 class="text-center text-7xl mt-3">Dont Touch My Dice!</h1>
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

<div class="fixed bottom-5 w-full text-center">Welcome to DTMD, the ultimate online platform for virtual dice games!<br>Create a room, invite a few of your friends and get rolling!</div>


<div class="card p-4 w-72 shadow-xl" data-popup="popupCreateRoom">
	<div class="flex flex-col">
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={roomname} placeholder="Room Name" />
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={usernameCreate} placeholder="Nickname" />
		{#if errorMessageCreate}
		<p class="text-red-500 mt-2">{errorMessageCreate}</p>
		{/if}
		<button class="btn variant-filled" style="margin-top: 2vh;" on:click={createRoom} disabled = {disableCreateButton} >Create</button>
	</div>
</div>

<div class="card p-4 w-72 shadow-xl" data-popup="popupJoinRoom">
	<div class="flex flex-col">
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={roomPin} placeholder="Room-Pin" />
		<input class="p-2.5 text-base m-2 rounded-full border-transparent text-black" type="text" bind:value={usernameJoin} placeholder="Nickname" />
		{#if errorMessageJoin}
		<p class="text-red-500 mt-2">{errorMessageJoin}</p>
		{/if}
		<button class="btn variant-filled" style="margin-top: 2vh;" on:click={joinRoom} disabled = {disableJoinButton} >Join</button>
	</div>
</div>
