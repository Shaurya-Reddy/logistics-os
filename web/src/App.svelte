<script lang="ts">
  import { onMount } from 'svelte';
  let status = $state<'loading' | 'ready' | 'unavailable'>('loading');
  async function check() {
    status = 'loading';
    try {
      const response = await fetch('/ready', { cache: 'no-store', signal: AbortSignal.timeout(4000) });
      status = response.ok ? 'ready' : 'unavailable';
    } catch { status = 'unavailable'; }
  }
  onMount(() => { void check(); });
</script>

<svelte:head><meta name="description" content="Logistics OS installation status" /></svelte:head>
<main>
  <header><span class="mark" aria-hidden="true">L</span><strong>Logistics OS</strong><span class="stage">S0 · Foundation</span></header>
  <section aria-labelledby="title">
    <p class="eyebrow">INSTALLATION STATUS</p>
    <h1 id="title">A foundation for<br />your operations.</h1>
    <p class="intro">Your Logistics OS application is running. This first milestone connects the application to its database.</p>
    <div class="status" aria-live="polite" aria-busy={status === 'loading'}>
      <span class:good={status === 'ready'} class="indicator" aria-hidden="true"></span>
      <div><h2>{status === 'ready' ? 'System ready' : status === 'loading' ? 'Checking connection…' : 'Database unavailable'}</h2>
        <p>{status === 'ready' ? 'Database connected and schema up to date.' : status === 'loading' ? 'Checking the database and installed schema.' : 'Ask your deployment administrator to check the database and migrations, then retry.'}</p></div>
    </div>
    <button onclick={check} disabled={status === 'loading'}>Check again <span aria-hidden="true">↗</span></button>
    <aside><h2>Next milestone</h2><p>Identity and master data: Organization → Site → Customer → Warehouse.</p><p>Business workflows are not available in this skeleton.</p></aside>
  </section>
  <footer>Self-hosted. Your data, your operations.</footer>
</main>

<style>
  :global(*){box-sizing:border-box} :global(body){margin:0;background:#f5f5f0;color:#173d36;font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}
  main{max-width:1080px;margin:auto;padding:32px 36px}header{display:flex;align-items:center;gap:12px;padding-bottom:28px;border-bottom:1px solid #d8dfd7}.mark{display:grid;place-items:center;background:#173d36;color:white;width:36px;height:36px;border-radius:8px;font-weight:700}.stage{margin-left:auto;font-size:13px;color:#52645b}section{max-width:650px;padding:68px 0 40px}.eyebrow{font-size:12px;letter-spacing:.16em;font-weight:700;color:#596e5e}h1{font-size:clamp(36px,6vw,60px);letter-spacing:-.045em;line-height:1.08;margin:20px 0}p{line-height:1.6}.intro{font-size:18px;color:#52645b;max-width:580px}.status{display:flex;gap:16px;background:white;border:1px solid #d8dfd7;border-radius:12px;padding:24px;margin:30px 0 16px}.indicator{flex-shrink:0;width:10px;height:10px;border-radius:50%;background:#9b722e;margin-top:8px}.indicator.good{background:#277453}h2{font-size:16px;margin:0}.status p{margin:5px 0 0;font-size:14px;color:#52645b}button{font:inherit;font-size:14px;background:#173d36;color:white;border:0;border-radius:7px;padding:12px 18px;cursor:pointer}button span{margin-left:20px}button:disabled{opacity:.65;cursor:wait}button:focus-visible{outline:3px solid #a66c18;outline-offset:4px}aside{margin-top:48px;border-top:1px solid #d8dfd7;padding-top:24px}aside p{font-size:14px;color:#52645b;margin:8px 0}footer{padding:24px 0;border-top:1px solid #d8dfd7;font-size:12px;color:#52645b}@media(max-width:540px){main{padding:20px}.stage{font-size:11px}section{padding-top:40px}}
</style>
