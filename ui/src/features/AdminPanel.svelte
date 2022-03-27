<script lang="ts">
  import DataTable, { Head, Body, Row, Cell } from "@smui/data-table";
  import { Link } from "svelte-navigator";
  import Loading from "../components/Loading.svelte";
  import type { Celebration } from "../types/Celebration";
  import { celebrations } from "../stores/Store";

  const fetchCelebrations = (async () => {
    const response = await fetch("http://localhost:8080/events");
    return (await response.json()) as Promise<Celebration[]>;
  })();

  function saveCelebrations(fetchedCelebrations: Celebration[]) {
    celebrations.set(fetchedCelebrations);
  }
</script>

<section>
  {#await fetchCelebrations}
    <Loading />
  {:then events}
    {saveCelebrations(events)}
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
        {#each events as event}
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
  {:catch error}
    <p>Server error while fetching the events</p>
  {/await}
</section>
