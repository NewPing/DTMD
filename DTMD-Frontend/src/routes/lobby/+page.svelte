<script lang="ts">
	import { LobbyID, MemberID } from './../../stores.js';
	import '../../app.postcss';
	import { SlideToggle, RadioGroup, RadioItem, type PopupSettings, popup } from '@skeletonlabs/skeleton';

	// Floating UI for Popups
	import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom';
	import { storePopup } from '@skeletonlabs/skeleton';
	import { Api, type MainRollDiceRequest,type MainChatMessage } from '../../dtmd_api';
	import { onMount } from 'svelte';
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
		placement: 'bottom',
	
	};
	//constants
	let updateMemberList:number = 0;
	let updateChat:number = 1;
	//style
	let lobbyName = '';
	let numberRolled = 42;
	let tempRollNumber : number;
	let isRolling = false;
	let receivedRoll = false;
	let chatMessages:MainChatMessage[] = [];
	let elemChat : HTMLElement;
	//lobbyState
	let members : string[] = [];
	let isPrivateMessage: boolean = false;
	//bind number variable to state of boolean variable cause backend needs a number
	$: isPrivateRoll = isPrivateMessage ? 1 : 0;
	let numberOfDice: number = 1;
	let diceType: number = 4;

	let lobbyID = '';
	let memberID = '';

	onMount(() => {
		//Load stored variables
		const unsubscribeLobbyID = LobbyID.subscribe(value => {https://discord.com/channels/@me/950465749590360064
			lobbyID = value;
		});

		const unsubscribeMemberID = MemberID.subscribe(value => {
			memberID = value;
		});
		unsubscribeLobbyID();
		unsubscribeMemberID();
		console.log("successfull join", lobbyID);
		//fetch lobby information
		fetchMembers().then( fetchedMembers => {
			members = fetchedMembers;
		});
		loadLobbyName().then( fetchedLobbyName => {
			lobbyName = fetchedLobbyName;
		});
		fetchUpdateRoutine();
	});

	async function fetchMembers(): Promise<string[]> {
		const res = await api.lobbiesMembersDetail(lobbyID);
		if (!res.ok) {
			throw new Error("Failed to fetch member list");
		}
		const response = await res.json();
		if (!Array.isArray(response)) {
			throw new Error('Error fetching members, not array of strings.');
		}
		return response; 
	}

	async function fetchChatMessages(): Promise<MainChatMessage[]> {
		const res = await api.lobbiesMembersMessagesDetail(lobbyID,memberID);
		if (!res.ok) {
			throw new Error("Failed to fetch chat messages.");
		}
		const response = await res.json();
		if (!Array.isArray(response)) {
			throw new Error('Error fetching members, not array of strings.');
		}
		return response; 
	}

	async function loadLobbyName(): Promise<string> {
		const res = await api.lobbiesNameDetail(lobbyID);
		if (!res.ok) {
			throw new Error("Failed to fetch lobby name");
		}
		const lobbyName = await res.text()
		return lobbyName;
	}
	function copyPinToClipboard(pin:string) {
		navigator.clipboard.writeText(pin)
 	}
	async function fetchUpdates(): Promise<number[]> {
		const res = await api.lobbiesMembersUpdatesDetail(lobbyID,memberID);
		if (!res.ok) {
			throw new Error("Failed to fetch member list");
		}
		const response = await res.json();
		if (!Array.isArray(response)) {
			throw new Error('Error fetching members, not array of strings.');
		}
		// Map and parse each item to ensure it's a number
		const parsedResponse = response.map(item => {
			// Assuming the items are strings that should be parsed to integers
			const parsedItem = parseInt(item, 10);
			if (isNaN(parsedItem)) {
				throw new Error('Data contains non-numeric values');
			}
			return parsedItem;
		});
		return parsedResponse;
	}
	async function postRoll() : Promise<number>{
		const rollDiceRequest : MainRollDiceRequest = {
			DiceType: diceType,
			IsPrivateRoll: isPrivateRoll,
			MemberID: memberID,
			NumberOfRolls: numberOfDice,
		};
		const res = await api.lobbiesRolldiceCreate(lobbyID,rollDiceRequest);
		if (!res.ok) {
			throw new Error("Failed to post number roll.");
		}
		const resultText = await res.text();
		const resultNumber = parseInt(resultText);
		return resultNumber;
	}

	function fetchUpdateRoutine(){
		//establish interval as timer function
		const interval = setInterval(() => {
			fetchUpdates().then( updates => {
				if(updates.includes(updateMemberList)){
					fetchMembers().then( fetchedMembers => {
					members = fetchedMembers;
					});
				}
				if(updates.includes(updateChat)){
					fetchChatMessages().then(fetchedChatMessages => {
						chatMessages = [...chatMessages,...fetchedChatMessages]
						// Timeout prevents race condition
						setTimeout(() => {
							scrollChatBottom('smooth');
						}, 0);
					});
				}
			});
		}, 500); 
			return () => clearInterval(interval);
	}

	function startRoll() {
		tempRollNumber = 0;
		receivedRoll = false;
		//make sure only one roll running at a time
		if (isRolling) return; 
		isRolling = true;
		var passedTime = 0;
		var changeInterval = 75;
		const rollTimer = setInterval(() => {
			passedTime += changeInterval;
			//display correct interval for possible results
			numberRolled = Math.floor(Math.random() * (numberOfDice * diceType - numberOfDice + 1)) + numberOfDice
			//show random numbers until we have the actual number from api
			if(receivedRoll){
				isRolling = false;
				numberRolled = tempRollNumber;
				clearInterval(rollTimer);
			}
			//send api request with 2 second delay to make it seem like actual roll is happening
			if(passedTime > 1500 && !receivedRoll)
			postRoll().then( rollResult => {
				console.log("roll result from backend",rollResult)
				//when receiving roll from API then notify that roll is received
				tempRollNumber = rollResult;
				receivedRoll = true;
			});
		}, changeInterval);
  	}
	function scrollChatBottom(behavior?: ScrollBehavior): void {
		elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}

</script>

<div class="grid h-screen grid-rows-[auto_1fr]">
	<!-- Header -->
	<header class="bg-surface-100-800-token p-4">
		<div class="grid grid-cols-2 gap-4">
		  <div class="font-bold text-xl">
			<span class="font-bold text-xl" style="color: #D4163C; font-weight: 900;">DTMD </span>
			<span class="font-bold text-xl">- {lobbyName === '' ? 'loading lobby name...' : lobbyName}</span>
		  </div>
		  <div class="text-right">
			<span class="">Lobby PIN:</span>
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<span class="hover:underline cursor-pointer" use:popup={popupClipboard} on:click={() => copyPinToClipboard(lobbyID)} style="color: #D4163C; font-weight: 700;">
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
				<button type="button" class="btn btn-lg variant-filled-primary font-semibold" on:click={startRoll} disabled = {isRolling}>Roll!</button>
				<!-- Big Number -->
				 <div style="margin-top: 100px">&nbsp;</div>
				<div class = "font-medium select-none" style="font-size: 240px;">
				{numberRolled}
				</div>
			</div>
		</main>

        <div bind:this={elemChat} class="bg-surface-500/5 p-4 overflow-y-auto col-span-2">
			{#each chatMessages as { message, sender }}
			<div class="card p-4 variant-soft mb-3">
			  <header class="flex justify-between items-center">
				<p class="font-bold">{sender}</p>
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
