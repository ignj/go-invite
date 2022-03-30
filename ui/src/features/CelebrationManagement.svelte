<script lang="ts">
  import DataTable, { Body, Cell, Head, Row } from "@smui/data-table";

  import { onMount } from "svelte";
  import Loading from "../components/Loading.svelte";
  import { celebrations } from "../stores/Store";
  import type { Celebration } from "../types/Celebration";

  let celebration: Celebration;
  const celebrationId = window.location.pathname.split("/").pop();

  onMount(async () => {
    celebrations.subscribe((c) => {
      console.log(c);
      return (celebration = c.find((e) => e._id === celebrationId));
    });

    if (!celebration) {
      const response = await fetch(
        "http://localhost:8080/events/celebrationId"
      );
      celebration = (await response.json()) as Celebration;
    }
    console.log(celebration);
    console.log(celebrations);
  });

  function transformAttendeeStatus(attendeeStatus: number): string {
    switch (attendeeStatus) {
      case 0:
        return "pending";
      case 1:
        return "confirmed";
      case 2:
        return "absent";
      default:
        return "unknown status";
    }
  }

  async function getInvitationLink(celebrationId: string, attendeeId: string) {
    const response = await fetch(
      `http://localhost:8080/events/${celebrationId}/invitation/${attendeeId}`
    );

    showTooltip(attendeeId);

    navigator.clipboard.writeText(
      `http://localhost:5000/invite/${(await response.json()).invitationHash}`
    );
  }

  function showTooltip(attendeeId: string) {
    var tooltip = document
      .getElementById(attendeeId)
      .getElementsByTagName("span")[0];

    tooltip.style.visibility = "visible";

    setTimeout(() => {
      tooltip.style.visibility = "hidden";
    }, 3000);
  }
</script>

<section>
  {#if celebration}
    <h1>Title: {celebration.title}</h1>
    <p>
      Description: {celebration.description}
    </p>
    <p>
      Date: {celebration.date}
    </p>
    <DataTable table$aria-label="Attendees" style="max-width: 100%;">
      <Head>
        <Row>
          <Cell>Id</Cell>
          <Cell>Full name</Cell>
          <Cell>Email</Cell>
          <Cell>Phone Number</Cell>
          <Cell>Status</Cell>
          <Cell />
        </Row>
      </Head>
      <Body>
        {#each celebration.attendees as attendee}
          <Row>
            <Cell>{attendee._id}</Cell>
            <Cell>{attendee.fullName}</Cell>
            <Cell>{attendee.email}</Cell>
            <Cell>{attendee.phoneNumber}</Cell>
            <Cell>{transformAttendeeStatus(attendee.attendStatus)}</Cell>
            <div class="tooltip" id={attendee._id}>
              <span class="tooltiptext">Copied to clipboard!</span>
              <Cell
                on:click={() => getInvitationLink(celebrationId, attendee._id)}
                >Get Invitation Link</Cell
              >
            </div>
          </Row>
        {/each}
      </Body>
    </DataTable>
  {:else}
    <Loading />
  {/if}
</section>

<style>
  /* Tooltip container */
  .tooltip {
    position: relative;
    display: inline-block;
    border-bottom: 1px dotted black; /* If you want dots under the hoverable text */
  }

  /* Tooltip text */
  .tooltip .tooltiptext {
    visibility: hidden;
    width: 120px;
    background-color: black;
    color: #fff;
    text-align: center;
    padding: 5px 0;
    border-radius: 6px;

    /* Position the tooltip text - see examples below! */
    position: absolute;
    z-index: 1;
  }

  /* Show the tooltip text when you mouse over the tooltip container */
  .tooltiptext {
    visibility: visible;
  }
</style>
