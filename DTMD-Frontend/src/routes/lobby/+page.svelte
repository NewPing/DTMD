<script lang="ts">
	import { LobbyID, MemberID, ErrorMessageStart } from './../../stores.js';
	import '../../app.postcss';
	import {
		SlideToggle,
		RadioGroup,
		RadioItem,
		type PopupSettings,
		popup
	} from '@skeletonlabs/skeleton';

	// Floating UI for Popups
	import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom';
	import { storePopup } from '@skeletonlabs/skeleton';
	import {
		Api,
		type MainRollDiceRequest,
		type ModelsChatMessage,
		type HttpResponse
	} from '../../dtmd_api';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });

	// A fallback for non-Vite environments
	const isServerSide = typeof window === 'undefined';

	const apiBase = new Api({
		baseUrl: isServerSide ? process.env.VITE_API_BASE_URL : import.meta.env.VITE_API_BASE_URL
	});

	const api = apiBase.api;

	const popupClipboard: PopupSettings = {
		event: 'click',
		target: 'popupClipboard',
		placement: 'bottom'
	};
	//constants
	let updateMemberList: number = 0;
	let updateChat: number = 1;
	//style
	let lobbyName = '';
	let numberRolled = 42;
	let tempRollNumber: number;
	let isRolling = false;
	let receivedRoll = false;
	let chatMessages: ModelsChatMessage[] = [];
	let elemChat: HTMLElement;
	//lobbyState
	let members: string[] = [];
	let isPrivateMessage: boolean = false;
	//bind number variable to state of boolean variable cause backend needs a number
	$: isPrivateRoll = isPrivateMessage ? 1 : 0;
	let numberOfDice: number = 1;
	let diceType: number = 4;

	let lobbyID = '';
	let memberID = '';

	onMount(() => {
		var initCheck = true;
		//Load stored variables
		const unsubscribeLobbyID = LobbyID.subscribe((value) => {
			if (value === null || value === undefined || value === '') {
				console.log('navigate back');
				goto('/');
				initCheck = false;
			}
			lobbyID = value;
		});

		const unsubscribeMemberID = MemberID.subscribe((value) => {
			if (value === null || value === undefined || value === '') {
				console.log('navigate back');
				goto('/');
				initCheck = false;
			}
			memberID = value;
		});
		unsubscribeLobbyID();
		unsubscribeMemberID();
		if (initCheck) {
			//fetch lobby information
			fetchMembers().then((fetchedMembers) => {
				members = fetchedMembers;
			});
			loadLobbyName().then((fetchedLobbyName) => {
				lobbyName = fetchedLobbyName;
			});
			fetchAllChatMessages().then((fetchedChatMessages) => {
				chatMessages = [...chatMessages, ...fetchedChatMessages];
			});
			fetchUpdateRoutine();
		}
	});

	async function fetchMembers(): Promise<string[]> {
		try {
			const res = await api.lobbiesMembersDetail(lobbyID);
			const response = await res.json();
			if (!Array.isArray(response)) {
				throw new Error('Error fetching members, not array of strings.');
			}
			return response;
		} catch (error) {
			throw error;
		}
	}

	async function fetchChatMessages(): Promise<ModelsChatMessage[]> {
		try {
			const res = await api.lobbiesMembersMessagesDetail(lobbyID, memberID);
			const response = await res.json();
			if (!Array.isArray(response)) {
				throw new Error('Error fetching chat messages, not an array.');
			}
			return response;
		} catch (error) {
			throw error;
		}
	}
	async function fetchAllChatMessages(): Promise<ModelsChatMessage[]> {
		try {
			const res = await api.lobbiesChathistoryDetail(lobbyID);
			const response = await res.json();
			if (!Array.isArray(response)) {
				throw new Error('Error fetching chat messages, not an array.');
			}
			return response;
		} catch (error) {
			throw error;
		}
	}

	async function loadLobbyName(): Promise<string> {
		try {
			const res = await api.lobbiesNameDetail(lobbyID);
			const lobbyName = await res.text();
			return lobbyName;
		} catch (error) {
			throw error;
		}
	}
	function copyPinToClipboard(pin: string) {
		navigator.clipboard.writeText(pin);
	}
	async function fetchUpdates(): Promise<number[]> {
		try {
			const res = await api.lobbiesMembersUpdatesDetail(lobbyID, memberID);
			const response = await res.json();
			if (!Array.isArray(response)) {
				throw new Error('Error fetching updates, not an array.');
			}
			// Map and parse each item to ensure it's a number
			const parsedResponse = response.map((item) => {
				// Assuming the items are strings that should be parsed to integers
				const parsedItem = parseInt(item, 10);
				if (isNaN(parsedItem)) {
					throw new Error('Data contains non-numeric values');
				}
				return parsedItem;
			});
			return parsedResponse;
		} catch (error) {
			throw error;
		}
	}
	async function postRoll(tmpDiceType: number, tmpNumberOfDice: number): Promise<number> {
		const rollDiceRequest: MainRollDiceRequest = {
			DiceType: tmpDiceType,
			IsPrivateRoll: isPrivateRoll,
			MemberID: memberID,
			NumberOfRolls: tmpNumberOfDice
		};
		try {
			const res = await api.lobbiesRolldiceCreate(lobbyID, rollDiceRequest);
			const resultText = await res.text();
			const resultNumber = parseInt(resultText);
			return resultNumber;
		} catch (error) {
			throw error;
		}
	}

	function fetchUpdateRoutine() {
		//establish interval as timer function
		const interval = setInterval(() => {
			fetchUpdates()
				.then((updates) => {
					if (updates.includes(updateMemberList)) {
						fetchMembers().then((fetchedMembers) => {
							members = fetchedMembers;
						});
					}
					if (updates.includes(updateChat)) {
						fetchChatMessages().then((fetchedChatMessages) => {
							chatMessages = [...chatMessages, ...fetchedChatMessages];
							// Timeout prevents race condition
							setTimeout(() => {
								scrollChatBottom('smooth');
							}, 0);
						});
					}
				})
				.catch((error) => {
					const errorResponse = error as HttpResponse<any, any>;
					console.log('error response fetch updates', errorResponse);
					if (errorResponse.status === 404 || errorResponse.status === undefined) {
						clearInterval(interval);
						LobbyID.set('');
						MemberID.set('');
						ErrorMessageStart.set('Room no longer exists.');
						goto('/');
					}
				});
		}, 500);
		return () => clearInterval(interval);
	}

	function startRoll(tmpDiceType: number, tmpNumberOfDice: number) {
		if (isRolling) return; // Ensure only one roll is running at a time
		isRolling = true;
		receivedRoll = false;
		tempRollNumber = 0;

		let passedTime = 0;
		let changeInterval = 35;
		var requestSent = false;
		function roll() {
			passedTime += changeInterval;
			// Generate a random number as a placeholder during the rolling animation
			numberRolled =
				Math.floor(Math.random() * (tmpNumberOfDice * tmpDiceType - tmpNumberOfDice + 1)) +
				tmpNumberOfDice;
			// Simulate delay to make it look like an actual roll is happening
			if (passedTime > 2000 && !receivedRoll && !requestSent) {
				requestSent = true;
				postRoll(tmpDiceType, tmpNumberOfDice).then((rollResult) => {
					console.log('Roll result from backend:', rollResult);
					tempRollNumber = rollResult; // Store the actual roll number
					receivedRoll = true; // Notify that the roll result is received
					// Immediately update numberRolled to show the final result
					numberRolled = tempRollNumber;
					isRolling = false; // Mark rolling as finished
				});
			}
			// Schedule the next roll with an increased interval if rolling is still ongoing
			if (passedTime < 2000 && !receivedRoll) {
				changeInterval *= 1.15; // Increase the interval to slow down the roll
				setTimeout(roll, changeInterval); // Schedule the next roll
			} else if (receivedRoll) {
				// Ensure the final number from the API is displayed
				numberRolled = tempRollNumber;
				isRolling = false;
			}
		}
		roll(); // Start the first roll
	}
	function scrollChatBottom(behavior?: ScrollBehavior): void {
		elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}
	function formatTime(timestamp?: string) {
		if (!timestamp) {
			return '';
		}
		const date = new Date(timestamp);
		return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
	}
</script>

<div class="grid h-screen grid-rows-[auto_1fr]">
	<!-- Header -->
	<header class="bg-surface-100-800-token p-4">
		<div class="grid grid-cols-2 gap-4">
			<div class="font-bold text-xl">
				<span class="font-bold text-xl" style="color: #D4163C; font-weight: 900;">DTMD </span>
				<span class="font-bold text-xl"
					>- {lobbyName === '' ? 'loading lobby name...' : lobbyName}</span
				>
			</div>
			<div class="text-right">
				<span class="">Lobby PIN:</span>
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				<span
					class="hover:underline cursor-pointer"
					use:popup={popupClipboard}
					on:click={() => copyPinToClipboard(lobbyID)}
					style="color: #D4163C; font-weight: 700;"
				>
					{lobbyID === '' ? 'loading Pin...' : lobbyID}
				</span>
			</div>
		</div>
	</header>
	<!-- Grid Columns -->
	<div class="grid grid-cols-6 overflow-hidden">
		<div class="bg-surface-500/5 p-4 overflow-y-auto">
			<nav class="list-nav" style="width: 100%;">
				{#if members.length === 0}
					<p>Loading members...</p>
				{:else}
					<ul>
						{#each members as item, _}
							<li>
								<!-- svelte-ignore a11y-missing-attribute --->
								<a class="flex-auto font-medium">
									<span>{item}</span>
								</a>
							</li>
						{/each}
					</ul>
				{/if}
			</nav>
		</div>

		<!-- Main Content -->
		<main class="space-y-4 p-4 flex flex-col overflow-auto col-span-3">
			<div class="flex flex-col items-center w-full pt-4">
				<SlideToggle name="slide" bind:checked={isPrivateMessage} active="bg-primary-500" size="sm"
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
				<div class="pt-7 mb-8">
					<h1 class="flex flex-col items-center mb-2">Dice Type</h1>
					<RadioGroup
						active="variant-filled-primary"
						hover="hover:variant-soft-primary"
						display="flex"
					>
						<RadioItem bind:group={diceType} name="justify" value={4}>4</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={6}>6</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={8}>8</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={10}>10</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={12}>12</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={16}>16</RadioItem>
						<RadioItem bind:group={diceType} name="justify" value={20}>20</RadioItem>
					</RadioGroup>
				</div>
				<button
					type="button"
					class="btn btn-lg variant-filled-primary font-semibold z-10"
					on:click={() => startRoll(diceType, numberOfDice)}
					disabled={isRolling}>Roll!</button
				>
				<!-- Big Number -->
				<div style="margin-top: 80px">&nbsp;</div>
				<div class="font-medium select-none z-0" style="font-size: 240px;">
					{numberRolled}
				</div>
			</div>
		</main>

		<div bind:this={elemChat} class="bg-surface-500/5 p-4 overflow-y-auto col-span-2">
			{#each chatMessages as { message, sender, timestamp }}
				<div class="card p-4 variant-soft mb-3">
					<header class="flex justify-between items-center">
						<p class="font-bold">{sender}</p>
						<small class="opacity-50">{formatTime(timestamp)}</small>
					</header>
					<p style="overflow-wrap: break-word;">{message}</p>
				</div>
			{/each}
		</div>
	</div>
</div>

<div class="card p-4 variant-filled-primary" data-popup="popupClipboard">
	<p>Copied to clipboard</p>
	<div class="arrow variant-filled-primary" />
</div>
