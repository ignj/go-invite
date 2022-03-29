<script lang="ts">
  const invitationHash = window.location.pathname.split("/").pop();

  async function confirmAttendance() {
    const response = await fetch(
      `http://localhost:8080/invitation/${invitationHash}/accept`,
      {
        method: "POST",
      }
    );
  }

  async function declineAttendance() {
    const response = await fetch(
      `http://localhost:8080/invitation/${invitationHash}/decline`,
      {
        method: "POST",
      }
    );
  }
</script>

<section>
  <h1>You have been invited to the ### event</h1>

  <h2>Please confirm your attendance</h2>
  <button on:click={confirmAttendance}>Confirm</button>
  <button on:click={declineAttendance}>Decline</button>
</section>
