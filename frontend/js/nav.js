import { auth, logoutUser } from "./firebase-config.js";
import { onAuthStateChanged } from "https://www.gstatic.com/firebasejs/9.23.0/firebase-auth.js";

function updateNav() {
  const navLinks = document.querySelectorAll(".nav-links");
  navLinks.forEach(nav => {
    // Remove existing auth links
    const existingAuthLinks = nav.querySelectorAll(".auth-nav-item");
    existingAuthLinks.forEach(link => link.remove());

    if (auth.currentUser) {
      // Show logout button
      const logoutLink = document.createElement("button");
      logoutLink.className = "button button-small button-secondary auth-nav-item";
      logoutLink.textContent = "Logout";
      logoutLink.type = "button";
      logoutLink.addEventListener("click", async () => {
        await logoutUser();
      });
      nav.appendChild(logoutLink);
    } else {
      // Show login and signup links
      const loginLink = document.createElement("a");
      loginLink.href = "auth.html";
      loginLink.className = "auth-nav-item";
      loginLink.textContent = "Login";

      const signupLink = document.createElement("a");
      signupLink.href = "auth.html";
      signupLink.className = "button button-small auth-nav-item";
      signupLink.textContent = "Signup";

      nav.appendChild(loginLink);
      nav.appendChild(signupLink);
    }
  });
}

// Initialize nav on page load
updateNav();
// Update nav when auth state changes
onAuthStateChanged(auth, () => {
  updateNav();
});
