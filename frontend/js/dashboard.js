import { API_BASE, authFetch, ensureAuth, handleAuthRejected, logoutUser } from "./firebase-config.js";

const dashboardPreview = document.getElementById("dashboard-preview");
const dashboardContent = document.getElementById("dashboard-content");
const bookingsList = document.getElementById("bookings-list");
const activityTimeline = document.getElementById("activity-timeline");
const miniCalendar = document.getElementById("mini-calendar");
const navAuthButton = document.getElementById("nav-auth-button");

function normalizeBookings(payload) {
  if (Array.isArray(payload)) return payload;
  if (Array.isArray(payload.bookings)) return payload.bookings;
  return [];
}

function bookingTitle(booking, index) {
  const slot = booking.slot || {};
  return booking.title || booking.service || booking.slotLabel || slot.time || `Booking ${index + 1}`;
}

function formatDate(dateStr) {
  if (!dateStr) return { day: "TBD", num: "—" };
  try {
    const date = new Date(dateStr);
    const days = ["SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"];
    return {
      day: days[date.getDay()],
      num: date.getDate().toString().padStart(2, "0")
    };
  } catch {
    return { day: "TBD", num: "—" };
  }
}

function renderBookings(bookings) {
  if (!bookingsList) return;
  bookingsList.innerHTML = "";

  if (!bookings.length) {
    const emptyDiv = document.createElement("div");
    emptyDiv.className = "dashboard-empty-state";
    emptyDiv.innerHTML = `
      <div class="dashboard-empty-state-icon">📅</div>
      <h3>No bookings yet</h3>
      <p>Book your first appointment to get started</p>
      <a class="button" href="slots.html">Browse Slots</a>
    `;
    bookingsList.appendChild(emptyDiv);
    return;
  }

  const fragment = document.createDocumentFragment();
  bookings.forEach((booking, index) => {
    const bookedSlot = booking.slot || {};
    const dateInfo = formatDate(booking.date || booking.day || bookedSlot.date);
    const time = bookedSlot.time || booking.slotId || booking.time || "Time TBD";

    const card = document.createElement("div");
    card.className = "upcoming-booking-card";
    card.innerHTML = `
      <div class="booking-left">
        <div class="booking-date-block">
          <p class="booking-date-day">${dateInfo.day}</p>
          <p class="booking-date-num">${dateInfo.num}</p>
        </div>
        <div class="booking-details">
          <p class="booking-time">${time}</p>
          <p class="booking-status">${bookingTitle(booking, index)}</p>
        </div>
      </div>
      <span class="booking-status-badge">${booking.status || "Confirmed"}</span>
    `;
    fragment.appendChild(card);
  });

  bookingsList.appendChild(fragment);
}

function renderActivityTimeline(bookings) {
  if (!activityTimeline) return;
  activityTimeline.innerHTML = "";

  const activities = [];
  
  if (bookings.length > 0) {
    activities.push({ text: "You booked a new appointment", time: "2 hours ago" });
  }
  
  activities.push({ text: "Account created", time: "Today" });

  const fragment = document.createDocumentFragment();
  activities.forEach(activity => {
    const item = document.createElement("div");
    item.className = "activity-item";
    item.innerHTML = `
      <div class="activity-dot"></div>
      <div class="activity-content">
        <p class="activity-text">${activity.text}</p>
        <p class="activity-time">${activity.time}</p>
      </div>
    `;
    fragment.appendChild(item);
  });

  activityTimeline.appendChild(fragment);
}

function renderMiniCalendar() {
  if (!miniCalendar) return;
  const now = new Date();
  const monthNames = ["January", "February", "March", "April", "May", "June",
                      "July", "August", "September", "October", "November", "December"];
  const month = now.getMonth();
  const year = now.getFullYear();
  const firstDay = new Date(year, month, 1).getDay();
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const today = now.getDate();

  let calendarHTML = `<div class="mini-calendar-title">${monthNames[month]} ${year}</div>`;
  calendarHTML += `<div class="mini-calendar-grid">`;
  
  const dayLabels = ["S", "M", "T", "W", "T", "F", "S"];
  dayLabels.forEach(day => {
    calendarHTML += `<div class="mini-calendar-day">${day}</div>`;
  });

  for (let i = 0; i < firstDay; i++) {
    calendarHTML += `<div class="mini-calendar-date" style="visibility: hidden;"></div>`;
  }

  const availableDates = [5, 8, 12, 15, 19, 22, 26];
  for (let i = 1; i <= daysInMonth; i++) {
    let classes = "mini-calendar-date";
    if (i === today) classes += " today";
    if (availableDates.includes(i)) classes += " available";
    calendarHTML += `<div class="${classes}">${i}</div>`;
  }

  calendarHTML += `</div>`;
  miniCalendar.innerHTML = calendarHTML;
}

function updateStats(bookings) {
  const totalBookingsEl = document.getElementById("total-bookings");
  const upcomingBookingsEl = document.getElementById("upcoming-bookings");
  const completedBookingsEl = document.getElementById("completed-bookings");

  if (totalBookingsEl) totalBookingsEl.textContent = bookings.length;
  if (upcomingBookingsEl) upcomingBookingsEl.textContent = bookings.filter(b => b.status !== "completed").length;
  if (completedBookingsEl) completedBookingsEl.textContent = bookings.filter(b => b.status === "completed").length;
}

async function fetchBookings() {
  try {
    const response = await authFetch(`${API_BASE}/api/bookings`);
    if (handleAuthRejected(response)) return;
    if (!response.ok) throw new Error(`Failed to load bookings`);
    const data = await response.json();
    const bookings = normalizeBookings(data);
    renderBookings(bookings);
    renderActivityTimeline(bookings);
    updateStats(bookings);
  } catch {
    // If backend is not available, show empty state
    renderBookings([]);
    renderActivityTimeline([]);
    updateStats([]);
  }
}

function showLoggedIn() {
  if (dashboardPreview) dashboardPreview.classList.add("hidden");
  if (dashboardContent) dashboardContent.classList.remove("hidden");
  
  if (navAuthButton) {
    navAuthButton.textContent = "Logout";
    navAuthButton.className = "button button-small button-secondary";
    navAuthButton.removeAttribute("href");
    navAuthButton.addEventListener("click", handleLogout);
  }
  
  renderMiniCalendar();
  fetchBookings();
}

function showLoggedOut() {
  if (dashboardContent) dashboardContent.classList.add("hidden");
  if (dashboardPreview) dashboardPreview.classList.remove("hidden");
  
  if (navAuthButton) {
    navAuthButton.textContent = "Sign In";
    navAuthButton.className = "button button-small";
    navAuthButton.setAttribute("href", "auth.html");
    navAuthButton.removeEventListener("click", handleLogout);
  }
}

async function handleLogout() {
  await logoutUser();
  showLoggedOut();
}

async function init() {
  let session = null;
  try {
    session = await ensureAuth();
  } catch {
    // Fail silently
  }

  if (session) {
    showLoggedIn();
  } else {
    showLoggedOut();
  }
}

init();
