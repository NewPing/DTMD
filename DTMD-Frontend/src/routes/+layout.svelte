<script lang="ts">
    import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
    import { storePopup } from '@skeletonlabs/skeleton';
	import Page from './+page.svelte';
    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

</script>

<Page></Page>