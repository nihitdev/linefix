import { useEffect, useState } from "react";

const REPOSITORY = "https://github.com/nihitdev/linefix";
const UNIX_INSTALL = "curl -fsSL https://raw.githubusercontent.com/nihitdev/linefix/main/install.sh | sh";
const WINDOWS_INSTALL = "irm https://raw.githubusercontent.com/nihitdev/linefix/main/install.ps1 | iex";

function Logo() {
  return <span className="brand-glyph" aria-hidden="true"><i /><i /></span>;
}

function Header({ docs = false }) {
  return (
    <header className="nav-wrap">
      <nav className="nav" aria-label="Primary navigation">
        <a className="brand" href={docs ? "/" : "#top"}><Logo /><span>linefix</span></a>
        <div className="nav-links">
          {docs ? <><a href="/#why">Why linefix</a><a href="/#install">Install</a><a href="/">Home</a></> : <><a href="#why">Why linefix</a><a href="#install">Install</a><a href="#commands">Commands</a><a href="/docs">Docs</a></>}
        </div>
        <a className="github-link" href={REPOSITORY} aria-label="View linefix on GitHub">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 .7a11.5 11.5 0 0 0-3.64 22.4c.58.1.79-.25.79-.56v-2.02c-3.22.7-3.9-1.37-3.9-1.37-.52-1.34-1.28-1.7-1.28-1.7-1.05-.72.08-.7.08-.7 1.16.08 1.77 1.19 1.77 1.19 1.03 1.77 2.7 1.26 3.36.96.1-.75.4-1.26.73-1.55-2.57-.29-5.27-1.28-5.27-5.68 0-1.26.45-2.28 1.19-3.09-.12-.29-.52-1.46.11-3.05 0 0 .97-.31 3.16 1.18a10.9 10.9 0 0 1 5.76 0c2.2-1.49 3.16-1.18 3.16-1.18.63 1.59.23 2.76.11 3.05.74.81 1.19 1.83 1.19 3.09 0 4.42-2.7 5.38-5.28 5.67.42.36.79 1.06.79 2.14v3.17c0 .31.21.67.8.56A11.5 11.5 0 0 0 12 .7Z" /></svg>
          <span>GitHub</span><b aria-hidden="true">↗</b>
        </a>
      </nav>
    </header>
  );
}

function Footer() {
  return (
    <footer>
      <a className="brand" href="/"><Logo /><span>linefix</span></a>
      <p>Line endings, fixed.</p>
      <div><a href="/docs">Docs</a><a href={`${REPOSITORY}/releases`}>Releases</a><a href={`${REPOSITORY}/blob/main/LICENSE`}>GPL-3.0-or-later</a></div>
    </footer>
  );
}

function Reveal({ as: Tag = "div", className = "", children, ...props }) {
  return <Tag className={`reveal ${className}`.trim()} {...props}>{children}</Tag>;
}

function FileCard({ after = false }) {
  return (
    <div className={`ending-card ${after ? "card-after" : "card-before"}`}>
      <span className="card-label">{after ? "after.txt" : "before.txt"}</span>
      <div className="text-lines"><span style={{ "--w": "78%" }} /><span style={{ "--w": "92%" }} /><span style={{ "--w": "64%" }} /><span style={{ "--w": "84%" }} /></div>
      <div className="ending-key">{after ? null : <i className="cr" />}<i className="lf" /><code>{after ? "LF" : "CRLF"}</code></div>
      {after ? <span className="safe-badge">✓ normalized</span> : null}
    </div>
  );
}

function Installer() {
  const [platform, setPlatform] = useState(() => /Win/i.test(navigator.userAgent) ? "windows" : "unix");
  const [copied, setCopied] = useState(false);
  const command = platform === "windows" ? WINDOWS_INSTALL : UNIX_INSTALL;

  async function copyCommand() {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
    } catch {
      setCopied(false);
    }
  }

  return (
    <Reveal className="installer delay-1">
      <div className="installer-tabs" role="tablist" aria-label="Installation platform">
        <button className={`install-tab ${platform === "unix" ? "active" : ""}`} type="button" role="tab" aria-selected={platform === "unix"} onClick={() => setPlatform("unix")}>Linux / macOS</button>
        <button className={`install-tab ${platform === "windows" ? "active" : ""}`} type="button" role="tab" aria-selected={platform === "windows"} onClick={() => setPlatform("windows")}>Windows</button>
      </div>
      <div className="install-panel" role="tabpanel">
        <div className="install-command"><span>{platform === "windows" ? ">" : "$"}</span><code>{command}</code><button className="copy" type="button" onClick={copyCommand} aria-label="Copy installation command"><span>{copied ? "Copied!" : "Copy"}</span></button></div>
        <div className="install-notes">{platform === "windows" ? <><span>✓ checksum verified</span><span>✓ per-user install</span><span>✓ PATH configured</span></> : <><span>✓ checksum verified</span><span>✓ no root required</span><span>✓ Intel + Apple Silicon</span></>}</div>
      </div>
      <a className="manual-download" href={`${REPOSITORY}/releases`}>Or download a release manually <b>↗</b></a>
    </Reveal>
  );
}

function LandingPage() {
  return (
    <>
      <Header />
      <main id="main">
        <section className="hero" id="top">
          <div className="hero-grid" aria-hidden="true" />
          <Reveal className="hero-content">
            <a className="release-pill" href={`${REPOSITORY}/releases/tag/v0.1.0`}><span>New</span> v0.1.0 is available <b>→</b></a>
            <h1>Every newline.<br /><em>Exactly right.</em></h1>
            <p className="hero-lede">Convert LF and CRLF without touching anything else. A focused, dependency-free CLI built for scripts, repositories, and humans.</p>
            <div className="hero-actions"><a className="button button-primary" href="#install">Install linefix <b>↓</b></a><a className="button button-ghost" href="/docs">Read the docs <b>→</b></a></div>
            <div className="hero-meta"><span><i className="pulse" /> Open source</span><span>Go standard library</span><span>Linux · macOS · Windows</span></div>
          </Reveal>
          <Reveal className="hero-visual delay-1" aria-label="linefix conversion example"><FileCard /><div className="conversion-arrow"><span>linefix lf</span><i>↓</i></div><FileCard after /></Reveal>
        </section>

        <section className="ticker" aria-label="linefix qualities"><div><span>NO CONFIG</span><i>✦</i><span>NO RUNTIME</span><i>✦</i><span>NO DEPENDENCIES</span><i>✦</i><span>NO SURPRISES</span><i>✦</i><span>NO CONFIG</span><i>✦</i><span>NO RUNTIME</span><i>✦</i></div></section>

        <section className="section why" id="why">
          <Reveal className="section-intro"><p className="eyebrow">Built to do one thing well</p><h2>Small surface.<br />Serious guarantees.</h2><p>linefix treats file contents with the care a low-level tool should. It changes the requested bytes and respects the rest.</p></Reveal>
          <div className="bento">
            <Reveal as="article" className="bento-card bento-main"><span className="card-icon">01</span><div className="binary-visual" aria-hidden="true"><span>01101100</span><span>01101001</span><span className="blocked">00000000</span><strong>binary?</strong><i>refuse</i></div><div><h3>Binary-safe by default</h3><p>Likely binary data is detected before modification. When the input looks wrong, linefix stops with a clear error.</p></div></Reveal>
            <Reveal as="article" className="bento-card delay-1"><span className="card-icon">02</span><div className="permission-visual" aria-hidden="true"><code>-rw-r-----</code><span>preserved</span></div><h3>Permissions stay put</h3><p>In-place conversion retains the file’s existing permission bits.</p></Reveal>
            <Reveal as="article" className="bento-card delay-2"><span className="card-icon">03</span><div className="skip-visual" aria-hidden="true"><span>LF</span><i>=</i><span>LF</span><strong>skip write</strong></div><h3>No pointless writes</h3><p>Already normalized? The file is left untouched—timestamps and all.</p></Reveal>
            <Reveal as="article" className="bento-card bento-wide"><span className="card-icon">04</span><div className="platforms"><span><b>⌁</b> Linux</span><span><b>⌘</b> macOS</span><span><b>⊞</b> Windows</span><span><b>64</b> AMD + ARM</span></div><div><h3>One binary. Every machine.</h3><p>Official AMD64 and ARM64 builds for all three major desktop platforms.</p></div></Reveal>
          </div>
        </section>

        <section className="install-section" id="install"><div className="install-inner"><Reveal className="section-intro light"><p className="eyebrow">Install in seconds</p><h2>Pick a shell.<br />Paste. Done.</h2><p>The installer identifies your system, downloads the correct release, and verifies its SHA-256 checksum before installing.</p></Reveal><Installer /></div></section>

        <section className="section commands" id="commands">
          <Reveal className="commands-heading"><p className="eyebrow">Command reference</p><h2>Three commands.<br />Zero guesswork.</h2></Reveal>
          <div className="command-stack"><Command index="01" syntax="linefix lf <file>" detail={<>CRLF <i>→</i> LF</>} /><Command index="02" syntax="linefix crlf <file>" detail={<>LF <i>→</i> CRLF</>} /><Command index="03" syntax="linefix check <file>" detail="Detect line endings" inspect /></div>
          <Reveal className="check-output"><span>check returns</span><code>LF</code><code>CRLF</code><code>Mixed</code><code>No line endings</code></Reveal>
        </section>

        <section className="final-cta"><div className="final-mark" aria-hidden="true"><span>CR</span><i /><span>LF</span></div><Reveal><p className="eyebrow">Your repository deserves clean diffs</p><h2>Make every newline intentional.</h2><p>Fast enough for scripts. Careful enough for source code.</p><div className="hero-actions centered"><a className="button button-primary" href="#install">Get linefix <b>↓</b></a><a className="button button-dark" href={REPOSITORY}>Star on GitHub <b>↗</b></a></div></Reveal></section>
      </main>
      <Footer />
    </>
  );
}

function Command({ index, syntax, detail, inspect = false }) {
  return <Reveal as="article" className="command-row"><span className="command-index">{index}</span><code><b>$</b> {syntax}</code><p>{detail}</p><span className={`command-tag ${inspect ? "inspect" : ""}`}>{inspect ? "inspect" : "convert"}</span></Reveal>;
}

const docsSections = [
  ["installation", "Installation"], ["commands-docs", "Command reference"], ["internals", "How it works"], ["development", "Development"], ["releasing", "Releasing"], ["troubleshooting", "Troubleshooting"],
];

function DocsPage() {
  return (
    <><Header docs /><main className="docs-shell"><aside className="docs-sidebar"><p>Documentation</p>{docsSections.map(([id, label]) => <a key={id} href={`#${id}`}>{label}</a>)}</aside><article className="docs-content"><p className="eyebrow">linefix docs</p><h1>Documentation</h1><p className="docs-lede">Everything needed to install, use, understand, and contribute to linefix.</p><div className="docs-callout"><strong>Quick start</strong><pre><code>linefix check README.md{"\n"}linefix lf notes.txt{"\n"}linefix crlf windows-file.txt</code></pre></div><DocsContent /></article></main><Footer /></>
  );
}

function DocsContent() {
  return <>
    <section id="installation"><h2>Installation</h2><h3>Linux and macOS</h3><pre><code>{UNIX_INSTALL}</code></pre><p>The checksum-verifying installer detects Linux or macOS and AMD64 or ARM64, then installs to <code>~/.local/bin</code>.</p><h3>Windows</h3><pre><code>{WINDOWS_INSTALL}</code></pre><p>Run in PowerShell. It installs per-user and adds the directory to PATH without Administrator privileges.</p></section>
    <section id="commands-docs"><h2>Command reference</h2><div className="docs-table"><div><code>linefix lf &lt;file&gt;</code><span>Convert CRLF to LF</span></div><div><code>linefix crlf &lt;file&gt;</code><span>Normalize and convert LF to CRLF</span></div><div><code>linefix check &lt;file&gt;</code><span>Print LF, CRLF, Mixed, or No line endings</span></div><div><code>linefix --version</code><span>Print the build version</span></div></div><h3>Exit codes</h3><ul><li><code>0</code> — success</li><li><code>1</code> — file or conversion error</li><li><code>2</code> — invalid arguments</li></ul></section>
    <section id="internals"><h2>How it works</h2><p>linefix reads a regular file, samples it for likely binary content, converts bytes in memory, and writes only when the result differs. Existing CRLF is normalized before CRLF conversion, so each newline is expanded exactly once.</p><h3>Safe replacement</h3><ol><li>Create a temporary file beside the original.</li><li>Apply the original permission bits.</li><li>Write and sync the converted content.</li><li>Replace the original path.</li></ol></section>
    <section id="development"><h2>Development</h2><p>Go 1.22 or newer is required. linefix has no third-party Go dependencies.</p><pre><code>git clone https://github.com/nihitdev/linefix.git{"\n"}cd linefix{"\n"}gofmt -w .{"\n"}go vet ./...{"\n"}go test ./...{"\n"}go build .</code></pre></section>
    <section id="releasing"><h2>Releasing</h2><p>Push an annotated semantic-version tag after CI passes:</p><pre><code>git tag -a v0.2.0 -m &quot;linefix v0.2.0&quot;{"\n"}git push origin v0.2.0</code></pre><p>The release workflow tests, cross-compiles all supported targets, packages each executable, and publishes checksums.</p></section>
    <section id="troubleshooting"><h2>Troubleshooting</h2><h3>Command not found</h3><p>Add <code>~/.local/bin</code> to PATH on Unix. On Windows, open a new terminal after installation.</p><h3>Binary-file rejection</h3><p>NUL bytes or a high proportion of control bytes trigger the safeguard. Convert legitimate text in unusual encodings to UTF-8 first.</p><h3>No line endings</h3><p>This is expected for empty files or files with no LF or CRLF sequence.</p></section>
  </>;
}

export default function App() {
  useEffect(() => {
    const elements = document.querySelectorAll(".reveal");
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add("visible");
          observer.unobserve(entry.target);
        }
      });
    }, { threshold: 0.1 });
    elements.forEach((element) => observer.observe(element));
    return () => observer.disconnect();
  }, []);

  return window.location.pathname.startsWith("/docs") ? <DocsPage /> : <LandingPage />;
}
