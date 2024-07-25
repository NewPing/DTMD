<script lang="ts">
	import '../app.postcss';
	import { popup } from '@skeletonlabs/skeleton';
	import type { PopupSettings } from '@skeletonlabs/skeleton';

	const popupFeatured: PopupSettings = {
		// Represents the type of event that opens/closed the popup
		event: 'click',
		// Matches the data-popup value on your popup element
		target: 'popupFeatured',
		// Defines which side of your trigger the popup will appear
		placement: 'bottom',
	};

	let roomname = '';
	let username = '';
</script>


<h1>Welcome to DTMD</h1>
<div class="centeredFrontPageContent">
	<button class="btn variant-filled createButton" style="margin-top: 10vh; font-size: 20px;" use:popup={popupFeatured}>Create Room</button>
	<div class="card p-4 w-72 shadow-xl" data-popup="popupFeatured">
	<div class="flex flex-col">
		<input type="text" bind:value={roomname} placeholder="Room Name">
		<input type="text" bind:value={username} placeholder="Nickname">
		<button class="btn variant-filled" style="margin-top: 2vh;">Create</button>
	</div>
	<div class="arrow bg-surface-100-800-token" />
</div>

</div>


<style lang="postcss">
	h1 {
		margin-top: 6vh;
		text-align: center;
		font-size: 40px;
	}
	.createButton {

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
    	justify-content: center;
    	align-items: center;
	}
	
</style>