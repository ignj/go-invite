<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table";
  import { Link } from "svelte-navigator";
  import Loading from "../components/Loading.svelte";
  import type { Celebration } from "../types/Celebration";
  import { celebrations } from "../stores/Store";
  import { onMount } from "svelte";

  let currentCelebrations;

  onMount(async () => {
    const response = await fetch("http://localhost:8080/events");
    celebrations.set((await response.json()) as Celebration[]);
  });

  celebrations.subscribe((value) => {
    currentCelebrations = value;
  });
</script>

<section>
  {#if currentCelebrations.length > 0}
    <h1>Celebrations</h1>
    <DataTable table$aria-label="Celebration list" style="max-width: 100%;">
      <Head>
        <Row>
          <Cell>Id</Cell>
          <Cell>Title</Cell>
          <Cell>Date</Cell>
          <Cell numeric>Attendees</Cell>
          <Cell />
        </Row>
      </Head>
      <Body>
        {#each currentCelebrations as event}
          <Row>
            <Cell>{event._id}</Cell>
            <Cell>{event.title}</Cell>
            <Cell>{event.date}</Cell>
            <Cell numeric>{event.attendees.length}</Cell>
            <Link to="/event/{event._id}">View more</Link>
          </Row>
        {/each}
      </Body>
    </DataTable>
  {:else}
    <Loading />
  {/if}
</section>
