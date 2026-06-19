import { ensureAuth } from "./firebase-config.js";

const authLoading = document.querySelector("#auth-loading");
const pageContent = document.querySelector("#page-content");

(async function init() {
  try {
    const session = await ensureAuth();

    if (session) {
      window.location.href = "dashboard.html";
      return;
    }
  } catch {
    // Keep the public landing page available if session validation fails here.
  }

  authLoading.classList.add("hidden");
  pageContent.classList.remove("hidden");
})();
