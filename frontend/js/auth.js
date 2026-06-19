const API_BASE = "http://localhost:8080";

import {
  auth,
  ensureAuth,
  loginWithEmail,
  loginWithGoogle,
  logoutUser,
  signupWithEmail,
  syncBackendUser
} from "./firebase-config.js";
import { onAuthStateChanged } from "https://www.gstatic.com/firebasejs/9.23.0/firebase-auth.js";

const tabs = document.querySelectorAll('[data-auth-mode]');
const authForm = document.querySelector('#auth-form');
const authSubmit = document.querySelector('#auth-submit');
const googleLoginBtn = document.querySelector('#google-login');
const logoutBtn = document.querySelector('#logout-button');
const statusMessage = document.querySelector('#auth-status');

let currentMode = 'login';

function setStatus(message, type = '') {
  statusMessage.textContent = message;
  statusMessage.className = `status-message ${type}`.trim();
}

tabs.forEach(tab => {
  tab.addEventListener('click', () => {
    currentMode = tab.dataset.authMode;
    
    tabs.forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    
    authSubmit.textContent = currentMode === 'login' ? 'Login' : 'Sign Up';
    setStatus('');
  });
});

authForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  
  const formData = new FormData(authForm);
  const email = formData.get('email');
  const password = formData.get('password');
  
  authSubmit.disabled = true;
  setStatus('');
  
  try {
    if (currentMode === 'login') {
      await loginWithEmail(email, password);
    } else {
      await signupWithEmail(email, password);
    }
    
    await syncBackendUser();
    
    setStatus('Successfully authenticated! Redirecting...', 'success');
    setTimeout(() => {
      window.location.href = 'dashboard.html';
    }, 1000);
  } catch (error) {
    setStatus(error.message, 'error');
  } finally {
    authSubmit.disabled = false;
  }
});

googleLoginBtn.addEventListener('click', async () => {
  googleLoginBtn.disabled = true;
  setStatus('');
  
  try {
    await loginWithGoogle();
    await syncBackendUser();
    
    setStatus('Successfully authenticated! Redirecting...', 'success');
    setTimeout(() => {
      window.location.href = 'dashboard.html';
    }, 1000);
  } catch (error) {
    setStatus(error.message, 'error');
  } finally {
    googleLoginBtn.disabled = false;
  }
});

logoutBtn.addEventListener('click', async () => {
  await logoutUser();
});

onAuthStateChanged(auth, (user) => {
  logoutBtn.classList.toggle('hidden', !user);
});

void API_BASE;
