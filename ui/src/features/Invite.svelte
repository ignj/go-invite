<script lang="ts">
  import { onMount } from "svelte";
  import { eventOverview } from "../stores/Store";
  import type { CelebrationOverview } from "../types/CelebrationOverview";

  const invitationHash = window.location.pathname.split("/").pop();
  let currentEvent: CelebrationOverview;
  let confirmedResponse: boolean = false;

  onMount(async () => {
    const response = await fetch(
      `http://localhost:8080/invitation/${invitationHash}/overview`
    );
    eventOverview.set((await response.json()) as CelebrationOverview);
  });

  eventOverview.subscribe((value) => {
    currentEvent = value;
  });

  async function confirmAttendance() {
    const response = await fetch(
      `http://localhost:8080/invitation/${invitationHash}/accept`,
      {
        method: "POST",
      }
    );

    if (response.status === 200) confirmedResponse = true;
  }

  async function declineAttendance() {
    const response = await fetch(
      `http://localhost:8080/invitation/${invitationHash}/decline`,
      {
        method: "POST",
      }
    );

    if (response.status === 200) confirmedResponse = true;
  }

  function showToast() {}
</script>

<section>
  {#if currentEvent}
    <h1>{currentEvent.title}</h1>
    <h2>You have been invited to attend on {currentEvent.date}!</h2>

    {#if confirmedResponse}
      <h3>Your response was saved</h3>
    {:else}
      <h3>Please confirm your attendance</h3>
      <button on:click={confirmAttendance}>Confirm</button>
      <button on:click={declineAttendance}>Decline</button>{/if}
  {/if}
</section>
