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
</script>


<h1>Welcome to DTMD</h1>
<div class="centeredFrontPageContent">
	<button class="btn variant-filled createButton" use:popup={popupFeatured}>Create Room</button>
	<div class="card p-4 w-72 shadow-xl" data-popup="popupFeatured">
	<div>
		<div>Test</div>
		<div>Test2</div>
		<button class="btn variant-filled" >TestButton</button>
		<div>test3</div>
		<p>Demo Content</p>
	</div>
	<div class="arrow bg-surface-100-800-token" />
</div>

</div>


<style lang="postcss">
	h1 {
		position: relative;
		text-align: center;
		font-size: 40px;
		//font-size: 3vw;
		top: 5vh;
	}
	.createButton{

	}
	.centeredFrontPageContent{
		height: 100vh;
    	display: flex;
    	justify-content: center;
    	align-items: center;
	}
	
</style>