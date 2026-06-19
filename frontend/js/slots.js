const API_BASE = "http://localhost:8080";

const slotsList = document.querySelector("#slots-list");
const refreshButton = document.querySelector("#refresh-slots");
const statusMessage = document.querySelector("#slots-status");

function setStatus(message, type = "") {
  statusMessage.textContent = message;
  statusMessage.className = `status-message ${type}`.trim();
}

function normalizeSlots(payload) {
  if (Array.isArray(payload)) {
    return payload;
  }

  if (Array.isArray(payload.slots)) {
    return payload.slots;
  }

  return [];
}

function slotId(slot) {
  return slot.id || slot.slotId || slot._id || slot.time || slot.label;
}

function slotTitle(slot) {
  return slot.label || slot.title || slot.time || `Slot ${slotId(slot) || ""}`.trim();
}

function slotDate(slot) {
  return slot.date || slot.day || slot.startDate || "Date to be confirmed";
}

function slotTime(slot) {
  return slot.time || slot.startTime || slot.range || "Time to be confirmed";
}

function renderSlots(slots) {
  slotsList.innerHTML = "";

  if (!slots.length) {
    const emptyState = document.createElement("div");
    const heading = document.createElement("h3");
    const copy = document.createElement("p");

    emptyState.className = "empty-state";
    heading.textContent = "No slots available";
    copy.textContent = "Try refreshing or check again later.";
    emptyState.append(heading, copy);
    slotsList.appendChild(emptyState);
    return;
  }

  const fragment = document.createDocumentFragment();

  slots.forEach((slot) => {
    const id = slotId(slot);
    const article = document.createElement("article");
    const header = document.createElement("div");
    const eyebrow = document.createElement("p");
    const title = document.createElement("h3");
    const meta = document.createElement("div");
    const date = document.createElement("p");
    const time = document.createElement("p");
    const link = document.createElement("a");

    article.className = "slot-card";
    eyebrow.className = "eyebrow";
    eyebrow.textContent = "Available";
    title.textContent = slotTitle(slot);
    header.append(eyebrow, title);

    meta.className = "slot-meta";
    date.append("Date: ", slotDate(slot));
    time.append("Time: ", slotTime(slot));
    meta.append(date, time);

    link.className = "button";
    link.href = `booking.html?slotId=${encodeURIComponent(id || "")}`;
    link.textContent = "Book";

    article.append(header, meta, link);
    fragment.appendChild(article);
  });

  slotsList.appendChild(fragment);
}

async function fetchSlots() {
  refreshButton.disabled = true;
  slotsList.innerHTML = "";
  setStatus("Loading available slots...");

  // Fallback slots for demo mode
  const fallbackSlots = [
    { id: 1, date: "2024-06-20", time: "10:00 AM" },
    { id: 2, date: "2024-06-20", time: "11:00 AM" },
    { id: 3, date: "2024-06-21", time: "2:00 PM" },
    { id: 4, date: "2024-06-22", time: "10:30 AM" },
    { id: 5, date: "2024-06-22", time: "1:30 PM" }
  ];

  try {
    const response = await fetch(`${API_BASE}/api/slots`);

    if (!response.ok) {
      throw new Error(`Could not load slots. Server returned ${response.status}.`);
    }

    const data = await response.json();
    let slots = normalizeSlots(data);
    
    // Use fallback slots if API returns empty
    if (!slots || slots.length === 0) {
      slots = fallbackSlots;
      setStatus("Showing demo slots (backend not connected).", "success");
    } else {
      setStatus(`${slots.length} slot${slots.length === 1 ? "" : "s"} loaded.`, "success");
    }
    
    renderSlots(slots);
  } catch (error) {
    // Use fallback slots on error
    renderSlots(fallbackSlots);
    setStatus("Showing demo slots (backend not connected).", "success");
  } finally {
    refreshButton.disabled = false;
  }
}

refreshButton.addEventListener("click", fetchSlots);
fetchSlots();
