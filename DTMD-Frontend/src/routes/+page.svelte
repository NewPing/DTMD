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

<label for="Name">Name:</label>

<h1>Welcome to DTMD</h1>

<style lang="postcss">
	h1 {
		position: relative;
		text-align: center;
		font-size: 3vw;
		top: 5vh;
	}
</style>

<button class="btn variant-filled" use:popup={popupFeatured}>Show Popup</button>
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