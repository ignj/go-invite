<script lang="ts">
  import { onMount } from "svelte";
  import { celebrations } from "../stores/Store";
  import type { Celebration } from "../types/Celebration";

  let celebration: Celebration = null;
  console.log("console.log(window.location.pathname.substring(1));");
  console.log(window.location.pathname.substring(1));

  celebrations.subscribe((c) => (celebration = c.find((e) => e._id == "")));
</script>

<h1>Celebration {celebration.title}</h1>
