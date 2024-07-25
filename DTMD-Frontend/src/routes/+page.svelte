<script lang="ts">
	import '../app.postcss';
	import { popup } from '@skeletonlabs/skeleton';
	import type { PopupSettings } from '@skeletonlabs/skeleton';

	const popupCreateRoom: PopupSettings = {
		// Represents the type of event that opens/closed the popup
		event: 'click',
		// Matches the data-popup value on your popup element
		target: 'popupCreateRoom',
		// Defines which side of your trigger the popup will appear
		placement: 'bottom'
	};
	const popupJoinRoom: PopupSettings = {
		// Represents the type of event that opens/closed the popup
		event: 'click',
		// Matches the data-popup value on your popup element
		target: 'popupJoinRoom',
		// Defines which side of your trigger the popup will appear
		placement: 'bottom'
	};

	let roomname = '';
	let username = '';
</script>

<h1>Welcome to DTMD</h1>
<div class="centeredFrontPageContent">
	<button
		class="btn variant-filled"
		style="margin-top: 30vh; font-size: 20px;"
		use:popup={popupCreateRoom}>Create Room</button
	>
	<button
	class="btn variant-filled"
	style="margin-top: 30px; font-size: 20px;"
	use:popup={popupJoinRoom}>Join Room</button>
</div>

<div class="card p-4 w-72 shadow-xl" data-popup="popupCreateRoom">
	<div class="flex flex-col">
		<input type="text" bind:value={roomname} placeholder="Room Name" />
		<input type="text" bind:value={username} placeholder="Nickname" />
		<button class="btn variant-filled" style="margin-top: 2vh;">Create</button>
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
	}
</style>
