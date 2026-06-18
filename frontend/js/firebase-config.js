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
  apiKey: "AIzaSyCrYg2tHmJR8eperOLgndykPJF-RXoZwto",
  authDomain: "booking-platform-943f9.firebaseapp.com",
  projectId: "booking-platform-943f9",
  storageBucket: "booking-platform-943f9.firebasestorage.app",
  messagingSenderId: "646294293293",
  appId: "1:646294293293:web:fc37b5af8d8a3521ada1a4",
  measurementId: "G-TZG35ZWDMB"
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

/**
 * Wait for Firebase Auth to finish restoring the persisted session from
 * IndexedDB.  The built-in `auth.authStateReady()` (Firebase v9.8+) returns
 * a promise that resolves **after** the initial auth state is determined.
 *
 * We fall back to a one-shot `onAuthStateChanged` listener for older SDK
 * builds where `authStateReady` may not exist.
 */
export function waitForAuthReady() {
  if (typeof auth.authStateReady === "function") {
    return auth.authStateReady().then(() => auth.currentUser);
  }

  // Fallback: wait for the first onAuthStateChanged callback.
  // NOTE: this is the *old* behaviour and still has the premature-null risk
  // if the SDK version is truly old, but it's the best we can do.
  return new Promise((resolve) => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      unsubscribe();
      resolve(user);
    });
  });
}

/**
 * Kept for backward compatibility with code that still calls it.
 * Now delegates to waitForAuthReady() so the premature-null race is fixed.
 */
export async function waitForAuthUser() {
  const user = await waitForAuthReady();

  if (user) {
    saveToken(await user.getIdToken());
  }

  return user;
}

/**
 * Master "restore session on page load" function.
 *
 * 1. Wait for Firebase Auth SDK to finish restoring from IndexedDB.
 * 2. If a currentUser exists → get a fresh token → return { user, token }.
 * 3. Otherwise, check localStorage for a previously-saved token and validate
 *    it against the backend's GET /api/me.
 * 4. If nothing works → return null (caller should redirect to auth.html).
 *
 * @returns {Promise<{user: object|null, token: string}|null>}
 */
export async function ensureAuth() {
  // Step 1 — wait for Firebase to finish its IndexedDB restoration.
  const firebaseUser = await waitForAuthReady();

  if (firebaseUser) {
    // Step 2 — Firebase recognised the persisted session.
    const token = await firebaseUser.getIdToken();
    saveToken(token);
    return { user: firebaseUser, token };
  }

  // Step 3 — Firebase didn't find a session.  Try the localStorage token.
  const storedToken = getStoredToken();
  if (!storedToken) {
    return null;
  }

  try {
    const response = await fetch(`${API_BASE}/api/me`, {
      headers: { Authorization: `Bearer ${storedToken}` }
    });

    if (response.ok) {
      // Token is still valid on the backend even though Firebase SDK lost
      // its client-side session (rare, but possible after browser data
      // partial-clear). Keep the token and let the caller proceed.
      return { user: null, token: storedToken };
    }

    // Token rejected (401 / expired / invalid) — clean up.
    clearToken();
    return null;
  } catch {
    // Network error — treat as unauthenticated rather than crashing.
    clearToken();
    return null;
  }
}

/**
 * Convenience wrapper around fetch() that attaches the Authorization header.
 * Uses getFreshIdToken() when a Firebase user is available, otherwise falls
 * back to the stored localStorage token.
 *
 * @param {string} url
 * @param {RequestInit} [options]
 * @returns {Promise<Response>}
 */
export async function authFetch(url, options = {}) {
  let token;

  try {
    token = await getFreshIdToken();
  } catch {
    token = getStoredToken();
  }

  if (!token) {
    throw new Error("No auth token available. Please sign in.");
  }

  const headers = new Headers(options.headers || {});
  headers.set("Authorization", `Bearer ${token}`);

  return fetch(url, { ...options, headers });
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

export async function syncBackendUser(user = auth.currentUser) {
	if (!user) {
		throw new Error("Please sign in before continuing.");
	}

	const token = await user.getIdToken();
	saveToken(token);

	const response = await fetch(`${API_BASE}/api/users/sync`, {
		method: "POST",
		headers: {
			Authorization: `Bearer ${token}`,
			"Content-Type": "application/json"
		},
		body: JSON.stringify({
			name: user.displayName || "",
			email: user.email || ""
		})
	});

	if (!response.ok) {
		let message = `Could not sync user. Server returned ${response.status}.`;

		try {
			const data = await response.json();
			message = data.error || message;
		} catch {
			// Keep the status-based message when the response body is not JSON.
		}

		throw new Error(message);
	}

	return response.json();
}

onAuthStateChanged(auth, async (user) => {
  if (user) {
    saveToken(await user.getIdToken());
    return;
  }

  clearToken();
});
