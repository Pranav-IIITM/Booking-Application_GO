const API_BASE = "http://localhost:8080";

import { ensureAuth, getFreshIdToken, logoutUser } from "./firebase-config.js";

const authLoading = document.querySelector("#auth-loading");
const pageContent = document.querySelector("#page-content");
const bookingsList = document.querySelector("#bookings-list");
const refreshButton = document.querySelector("#refresh-bookings");
const logoutButton = document.querySelector("#logout-button");
const statusMessage = document.querySelector("#dashboard-status");

function setStatus(message, type = "") {
  statusMessage.textContent = message;
  statusMessage.className = `status-message ${type}`.trim();
}

function normalizeBookings(payload) {
  if (Array.isArray(payload)) {
    return payload;
  }

  if (Array.isArray(payload.bookings)) {
    return payload.bookings;
  }

  return [];
}

function bookingTitle(booking, index) {
  const slot = booking.slot || {};
  return booking.title || booking.service || booking.slotLabel || slot.time || `Booking ${index + 1}`;
}

function renderBookings(bookings) {
  bookingsList.innerHTML = "";

  if (!bookings.length) {
    const emptyState = document.createElement("div");
    const heading = document.createElement("h3");
    const copy = document.createElement("p");

    emptyState.className = "empty-state";
    heading.textContent = "No bookings yet";
    copy.textContent = "Your confirmed reservations will appear here.";
    emptyState.append(heading, copy);
    bookingsList.appendChild(emptyState);
    return;
  }

  const fragment = document.createDocumentFragment();

  bookings.forEach((booking, index) => {
    const bookedSlot = booking.slot || {};
    const article = document.createElement("article");
    const header = document.createElement("div");
    const eyebrow = document.createElement("p");
    const title = document.createElement("h3");
    const meta = document.createElement("div");
    const name = document.createElement("p");
    const date = document.createElement("p");
    const slot = document.createElement("p");

    article.className = "booking-card";
    eyebrow.className = "eyebrow";
    eyebrow.textContent = "Confirmed";
    title.textContent = bookingTitle(booking, index);
    header.append(eyebrow, title);

    meta.className = "booking-meta";
    name.append("Status: ", booking.status || "confirmed");
    date.append("Date: ", booking.date || booking.day || bookedSlot.date || "Date to be confirmed");
    slot.append("Slot: ", bookedSlot.time || booking.slotId || booking.time || "Slot to be confirmed");
    meta.append(name, date, slot);

    article.append(header, meta);
    fragment.appendChild(article);
  });

  bookingsList.appendChild(fragment);
}

async function fetchBookings() {
  refreshButton.disabled = true;
  bookingsList.innerHTML = "";
  setStatus("Loading your bookings...");

  try {
    const token = await getFreshIdToken();
    const response = await fetch(`${API_BASE}/api/bookings`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });

    if (response.status === 401) {
      window.location.href = "auth.html";
      return;
    }

    if (!response.ok) {
      throw new Error(`Could not load bookings. Server returned ${response.status}.`);
    }

    const data = await response.json();
    const bookings = normalizeBookings(data);
    renderBookings(bookings);
    setStatus(`${bookings.length} booking${bookings.length === 1 ? "" : "s"} loaded.`, "success");
  } catch (error) {
    renderBookings([]);
    setStatus(error.message, "error");
  } finally {
    refreshButton.disabled = false;
  }
}

refreshButton.addEventListener("click", fetchBookings);
logoutButton.addEventListener("click", async () => {
  await logoutUser();
  window.location.href = "auth.html";
});

// ── Auth gate: restore session before loading data ──────────────────────
(async function init() {
  const session = await ensureAuth();

  if (!session) {
    window.location.href = "auth.html";
    return;
  }

  // Auth confirmed — reveal the page and fetch data.
  authLoading.classList.add("hidden");
  pageContent.classList.remove("hidden");
  fetchBookings();
})();
