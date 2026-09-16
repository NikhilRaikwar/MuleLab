import Link from "next/link";

export default function Home() {
  return <main>
    <section className="hero">
      <div className="eyebrow">OUTCOME-ACCOUNTABLE AUTONOMY</div>
      <h1>Give your business a goal.<br/><em>Make AI agents prove they deserve to stay.</em></h1>
      <p>MuleLab is a deterministic commerce sandbox where specialist agents propose measurable experiments, use bounded tools, and keep their authority only when the numbers support it.</p>
      <div className="actions"><Link className="button primary" href="/demo">Run the 2-minute demo →</Link><Link className="button" href="/architecture">Inspect the system</Link></div>
    </section>
    <section className="proofGrid">
      <article><b>01</b><h2>Models propose</h2><p>Typed hypotheses, evidence, budgets, tools, thresholds, and stop conditions.</p></article>
      <article><b>02</b><h2>Code authorizes</h2><p>Policy, permissions, approvals, idempotency, and spend limits are deterministic.</p></article>
      <article><b>03</b><h2>Outcomes decide</h2><p>A seeded simulator measures lift and retires strategies that fail their KPI.</p></article>
    </section>
  </main>;
}

