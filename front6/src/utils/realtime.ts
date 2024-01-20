export function getEvents(viewerId: string, listener: (data: string) => void) {
    let basePath = '';
    if (window.location.hostname === "localhost") {
        basePath = 'http://localhost:8001'
    } else {
        basePath = 'https://meme.mmaks.me/realtime';
    }

    let eventSource = new EventSource(basePath +"/listen?key=" + viewerId);

    eventSource.onmessage = function(event) {
        const parsed = JSON.parse(event.data);
        listener(parsed);
    };
}
