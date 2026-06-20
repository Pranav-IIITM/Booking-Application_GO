# Booking-Application_GO
A full-stack booking platform

# Roles

- Pranav : Backend (DB + Middlewares + *logic*)
- Barbie : Front-end + Backend(*Logic*)

---

# to end the previous server 

$pid = (netstat -ano | findstr :8080 | Select-Object -First 1) -split '\s+' | Select-Object -Last 1; if ($pid) { taskkill /PID $pid /F }; go run .