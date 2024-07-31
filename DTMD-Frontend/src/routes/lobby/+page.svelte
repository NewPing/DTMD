<script lang="ts">
	import { LobbyID, MemberID } from './../../stores.js';
	import '../../app.postcss';
	import { SlideToggle, RadioGroup, RadioItem } from '@skeletonlabs/skeleton';

	// Floating UI for Popups
	import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom';
	import { storePopup } from '@skeletonlabs/skeleton';
	import { Api, /*type MainMember*/ } from '../../dtmd_api';
	import { onMount } from 'svelte';
	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });

	const api = new Api({
		baseUrl: "http://localhost:8080",
	});

	let isPrivateMessage: boolean = false;
	let numberOfDice: number = 1;
	let diceType: number = 1;

	let lobbyID = '';
	let memberID = '';

	onMount(() => {
		const unsubscribeLobbyID = LobbyID.subscribe(value => {
			lobbyID = value;
		});

		const unsubscribeMemberID = MemberID.subscribe(value => {
			memberID = value;
		});

		unsubscribeLobbyID();
		unsubscribeMemberID();

		console.log("successfull join", lobbyID);
	});

	/*async function loadMembers(): Promise<MainMember[]> {
		const res = await api.members.membersList();
		if (res.ok) {
			return res.data;
			
		} else {
			throw new Error("Failed to fetch member list");
		}
	}*/
</script>


<div class="grid h-screen grid-rows-[auto_1fr_auto]">
	<!-- Header -->
	<header class="bg-surface-100-800-token p-4">
		<div class="grid grid-cols-2 gap-4">
			<div class="font-bold text-xl">DTMD - My Lobby 123</div>
			<div class="text-right">Lobby PIN: 123456</div>
		</div>
	</header>
	<!-- Grid Columns -->
	<div class="grid grid-cols-6">
		<div class="bg-surface-500/5 p-4 flex h-screen overflow-y-auto">
			<nav class="list-nav" style="width: 100%;">
				
			<!--	{#await loadMembers()}
					<p>Loading members...</p>
				{:then members}
					<ul>
						{#each members as item, i}
							<li>
								<!-- svelte-ignore a11y-missing-attribute
								<a class="flex-auto font-medium">
									<span>{item.name}</span>
								</a>
							</li>
						{/each}
					</ul>
				{:catch error}
					<p style="color: red">Error: {error.message}</p>
				{/await } -->
			</nav>
		</div>

		<!-- Main Content -->
		<main class="space-y-4 p-4 flex h-screen overflow-auto col-span-3">
			<div class="flex flex-col items-center w-full pt-4">
				<SlideToggle name="slide" bind:checked={isPrivateMessage} size="sm"
					>{isPrivateMessage ? 'Private' : 'Public'} Roll</SlideToggle
				>
				<div class="pt-7">
					<h1 class="flex flex-col items-center">Number of Rolls</h1>
					<RadioGroup active="variant-filled-primary" hover="hover:variant-soft-primary">
						<RadioItem bind:group={numberOfDice} name="justify" value={1}>1</RadioItem>
						<RadioItem bind:group={numberOfDice} name="justify" value={2}>2</RadioItem>
						<RadioItem bind:group={numberOfDice} name="justify" value={3}>3</RadioItem>
                        <RadioItem bind:group={numberOfDice} name="justify" value={4}>4</RadioItem>
                        <RadioItem bind:group={numberOfDice} name="justify" value={5}>5</RadioItem>
					</RadioGroup>
				</div>

				<div class="pt-7">
					<h1 class="flex flex-col items-center mb-2">Dice Type</h1>
					<RadioGroup
						active="variant-filled-primary"
						hover="hover:variant-soft-primary"
						display="flex"
					>
						<RadioItem bind:group={diceType} name="justify" value={1}>4</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={2}>6</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={3}>8</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={4}>10</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={5}>12</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={6}>16</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={7}>20</RadioItem>
					</RadioGroup>
				</div>
				<!-- Big Number -->
				 <div style="margin-top: 130px">&nbsp;</div>
				<div class = "font-medium select-none" style="font-size: 240px;">
					42
				</div>
			</div>
		</main>

		<div class="bg-surface-500/5 p-4 h-screen overflow-y-auto col-span-2">
			<div class="card p-4 variant-soft mb-3">
				<header class="flex justify-between items-center">
					<p class="font-bold" style="color: lightblue">AromaticA</p>
				</header>
				<p style="overflow-wrap: break-word;">Hello World!</p>
			</div>

			<div class="card p-4 variant-soft mb-3">
				<header class="flex justify-between items-center">
					<p class="font-bold" style="color: lightgreen">JullyJ</p>
				</header>
				<p style="overflow-wrap: break-word;">Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s.</p>
			</div>

			<div class="card p-4 variant-soft mb-3">
				<header class="flex justify-between items-center">
					<p class="font-bold" style="color: lightpink">PeterBliat</p>
				</header>
				<p style="overflow-wrap: break-word;">Where does it come from?</p>
			</div>

			<div class="card p-4 variant-soft mb-3">
				<header class="flex justify-between items-center">
					<p class="font-bold" style="color: lightblue">AromaticA</p>
				</header>
				<p style="overflow-wrap: break-word;">Contrary to popular belief, Lorem Ipsum is not simply random text.</p>
			</div>
			
		</div>
	</div>
	<!-- Footer -->
	<footer class="bg-blue-500 p-4">(footer)</footer>
</div>