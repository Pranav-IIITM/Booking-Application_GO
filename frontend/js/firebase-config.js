const API_BASE = "http://localhost:8080";

import { initializeApp } from "https://www.gstatic.com/firebasejs/9.23.0/firebase-app.js";
import {
  GoogleAuthProvider,
  createUserWithEmailAndPassword,
  getAuth,
  onAuthStateChanged,
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut
} from "https://www.gstatic.com/firebasejs/9.23.0/firebase-auth.js";

const firebaseConfig = {
  apiKey: "YOUR_FIREBASE_API_KEY",
  authDomain: "YOUR_FIREBASE_AUTH_DOMAIN",
  projectId: "YOUR_FIREBASE_PROJECT_ID",
  appId: "YOUR_FIREBASE_APP_ID"
};

const TOKEN_KEY = "firebaseIdToken";

export { API_BASE };

export const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
export const googleProvider = new GoogleAuthProvider();

export function saveToken(token) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export function getStoredToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export async function getFreshIdToken() {
  if (!auth.currentUser) {
    throw new Error("Please sign in before continuing.");
  }

  const token = await auth.currentUser.getIdToken();
  saveToken(token);
  return token;
}

export function waitForAuthUser() {
  return new Promise((resolve) => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      unsubscribe();

      if (user) {
        saveToken(await user.getIdToken());
      }

      resolve(user);
    });
  });
}

export async function loginWithEmail(email, password) {
  const credential = await signInWithEmailAndPassword(auth, email, password);
  saveToken(await credential.user.getIdToken());
  return credential.user;
}

export async function signupWithEmail(email, password) {
  const credential = await createUserWithEmailAndPassword(auth, email, password);
  saveToken(await credential.user.getIdToken());
  return credential.user;
}

export async function loginWithGoogle() {
  const credential = await signInWithPopup(auth, googleProvider);
  saveToken(await credential.user.getIdToken());
  return credential.user;
}

export async function logoutUser() {
  await signOut(auth);
  clearToken();
}

onAuthStateChanged(auth, async (user) => {
  if (user) {
    saveToken(await user.getIdToken());
    return;
  }

  clearToken();
});
