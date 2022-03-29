<script lang="ts">
  import { Router, Route } from "svelte-navigator";
  import AdminPanel from "./features/AdminPanel.svelte";
  import CelebrationManagement from "./features/CelebrationManagement.svelte";
  import Invite from "./features/Invite.svelte";
  import NotFound from "./components/NotFound.svelte";
</script>

<main>
  <Router>
    <Route path="/admin" component={AdminPanel} />
    <Route path="/event/:id" component={CelebrationManagement} />
    <Route path="/invite/:hash" component={Invite} />
    <Route path="/*" component={NotFound} />
  </Router>
</main>

<style>
  main {
    display: flex;
    padding: 1em;
    justify-content: center;
  }
</style>
