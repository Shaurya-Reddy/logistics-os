<script lang="ts">
  import { onMount } from 'svelte';

  let health = 'Checking…';
  let readiness = 'Checking…';

  async function readStatus(path: string) {
    const response = await fetch(path, { headers: { Accept: 'application/json' } });
    const body = await response.json();
    return response.ok ? body.status : 'not_ready';
  }

  onMount(async () => {
    try {
      health = await readStatus('/health');
    } catch {
      health = 'unavailable';
    }

    try {
      readiness = await readStatus('/ready');
    } catch {
      readiness = 'not_ready';
    }
  });
</script>

<svelte:head>
  <meta name="description" content="Logistics OS system shell" />
</svelte:head>

<main class="shell">
  <header>
    <div>
      <p class="eyebrow">Logistics OS</p>
      <h1>System foundation</h1>
      <p class="subtitle">A lightweight, self-hosted logistics operating system.</p>
    </div>
    <span class="version">v0.0.1-dev</span>
  </header>

  <section class="status-grid" aria-label="System status">
    <article>
      <span>Application</span>
      <strong>{health}</strong>
    </article>
    <article>
      <span>Database</span>
      <strong>{readiness === 'ready' ? 'connected' : readiness}</strong>
    </article>
    <article>
      <span>Runtime</span>
      <strong>Go + PostgreSQL</strong>
    </article>
  </section>

  <section class="next">
    <h2>S0 technical skeleton</h2>
    <p>Business modules are intentionally not enabled yet. The next gate is to prove startup, database readiness, migrations and performance budgets.</p>
  </section>
</main>
