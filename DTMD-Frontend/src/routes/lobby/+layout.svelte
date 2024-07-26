<script lang="ts">
	import '../../app.postcss';
	import { AppShell, AppBar } from '@skeletonlabs/skeleton';

	// Floating UI for Popups
	import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom';
	import { storePopup } from '@skeletonlabs/skeleton';
	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });
</script>


<div class="grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Header -->
	<header class="bg-surface-100-800-token p-4">
		<div class = "grid grid-cols-2 gap-4">
			<div class = "font-bold text-xl">LobbyName</div>
			<div class = "text-right ">LobbyPin</div>
		</div>
	</header>
	<!-- Grid Columns -->
	<div class="grid grid-cols-1 md:grid-cols-[auto_1fr]">
	  <!-- Left Sidebar. -->
	  <aside class="bg-surface-500/5 p-4 ">
		<nav class="list-nav">
			<!-- (optionally you can provide a label here) -->
			<ul>
				<li>
					<!-- svelte-ignore a11y-missing-attribute -->
					<a class="flex-auto font-medium">
						<span>Member</span>
					</a>
				</li>
				<!-- ... -->
			</ul>
		</nav>
	  </aside>
	  <!-- Main Content -->
	  <main class="space-y-4  p-4">
		<p class="h-[512px]  p-4">Paragraph 1</p>
		<p class="h-[512px]  p-4">Paragraph 2</p>
		<p class="h-[512px]  p-4">Paragraph 3</p>
	  </main>
	</div>
	<!-- Footer -->
	<footer class="bg-blue-500 p-4">(footer)</footer>
  </div>
  

