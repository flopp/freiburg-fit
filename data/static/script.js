document.addEventListener('DOMContentLoaded', function() {
    // MAPS
    var mapDiv = document.getElementById('venue-map');
    if (mapDiv) {
        var lat = parseFloat(mapDiv.dataset.lat);
        var lon = parseFloat(mapDiv.dataset.lon);
        var name = mapDiv.dataset.name;
        var map = L.map('venue-map', {gestureHandling: true}).setView([lat, lon], 13);
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            maxZoom: 18,
            attribution: '© OpenStreetMap contributors'
        }).addTo(map);
        L.marker([lat, lon]).addTo(map)
            .bindPopup(name).openPopup();
    }

    mapDiv = document.getElementById('venues-map');
    if (mapDiv) {
        // collect venue data (all html elements with data-venue)
        var venueData = [];
        document.querySelectorAll('[data-lat]').forEach(function(venueEl) {
            // skip elements without coordinates
            if (!venueEl.dataset.lat || !venueEl.dataset.lon) return;
            venueData.push({
                url: venueEl.dataset.url,
                name: venueEl.dataset.name,
                clubs: venueEl.dataset.clubs,
                lat: venueEl.dataset.lat,
                lon: venueEl.dataset.lon
            });
        });
        const freiburg = [47.996090, 7.849400];
        var lls = [freiburg];
        venueData.forEach(function(venue) {
            lls.push([venue.lat, venue.lon]);
        });
        var bounds = L.latLngBounds(lls).pad(0.3);
        var map = L.map('venues-map', {gestureHandling: true}).fitBounds(bounds);
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            maxZoom: 18,
            attribution: '© OpenStreetMap contributors'
        }).addTo(map);
        venueData.forEach(function(venue) {
            L.marker([venue.lat, venue.lon]).addTo(map)
                .bindPopup('<a href="' + venue.url + '">' + venue.name + '</a>');
        });
    }

    // UMAMI
    document.querySelectorAll("a[target=_blank]").forEach((a) => {
        if (a.getAttribute("data-umami-event") === null) {
            a.setAttribute('data-umami-event', 'outbound-link-click');
        }
        a.setAttribute('data-umami-event-url', a.href);
    });
    if (location.hash === '#disable-umami') {
        localStorage.setItem('umami.disabled', 'true');
        alert('Umami is now DISABLED in this browser.');
    }
    if (location.hash === '#enable-umami') {
        localStorage.removeItem('umami.disabled');
        alert('Umami is now ENABLED in this browser.');
    }
});
