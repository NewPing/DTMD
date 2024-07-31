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
		baseUrl: "http://localhost:8080",
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
	let userID = '';
	let roomname = '';
	let username = '';
	let errorMessage = '';
	let disableCreateButton = false;
	async function createRoom() {
		//Check input
		if(roomname === ''){
			errorMessage = "Please enter a room name.";
			setTimeout(() => {
                errorMessage = '';
            }, 3000);
			return;
		}
		if(username === ''){
			errorMessage = "Please enter a nickname.";
			setTimeout(() => {
                errorMessage = '';
            }, 3000);
			return;
		}
		disableCreateButton=true;
		var pin = await createLobbyAPICall(roomname);
		var userID = await joinLobbyAPICall(pin,username);
		MemberID.set(userID);
		LobbyID.set(pin);
		console.log("successfull creation",userID);
		disableCreateButton=false;
		goto('/lobby');
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

<h1>Dont Touch My Dice!</h1>
<div class="centeredFrontPageContent">
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

<div class="footer" >Welcome to DTMD, the ultimate online platform for virtual dice games!<br>Create a room, invite a few of your friends and get rolling!</div>


<div class="card p-4 w-72 shadow-xl" data-popup="popupCreateRoom">
	<div class="flex flex-col">
		<input type="text" bind:value={roomname} placeholder="Room Name" />
		<input type="text" bind:value={username} placeholder="Nickname" />
		{#if errorMessage}
		<p class="text-red-500 mt-2">{errorMessage}</p>
		{/if}
		<button class="btn variant-filled" style="margin-top: 2vh;" on:click={createRoom} disabled = {disableCreateButton} >Create</button>
	</div>
</div>

<div class="card p-4 w-72 shadow-xl" data-popup="popupJoinRoom">
	<div class="flex flex-col">
		<input type="text" bind:value={roomname} placeholder="Room-Pin" />
		<input type="text" bind:value={username} placeholder="Nickname" />
		<button class="btn variant-filled" style="margin-top: 2vh;">Join</button>
	</div>
</div>

<style lang="postcss">
	h1 {
		margin-top: 3vh;
		text-align: center;
		font-size: 60px;
		line-height: 60px;
	}

	input {
		padding: 10px;
		font-size: 16px;
		margin: 10px;
		border-radius: 25px;
		border-color: transparent;
		color: black;
	}

	.centeredFrontPageContent {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		text-align: center;
	}

	.footer {
		position: fixed;
		bottom: 0;
		margin-bottom: 3vh;
		width: 100%;
		text-align: center;
	}
</style>
