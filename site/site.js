document.querySelectorAll("[data-copy]").forEach((block) => {
  const button = block.querySelector(".copy");
  button.addEventListener("click", async () => {
    await navigator.clipboard.writeText(block.dataset.copy);
    button.textContent = "Copied";
    setTimeout(() => { button.textContent = "Copy"; }, 1600);
  });
});
