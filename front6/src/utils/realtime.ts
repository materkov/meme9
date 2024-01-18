export function getEvents(viewerId: string, listener: (data: string) => void) {
    //let eventSource = new EventSource("http://localhost:8001/listen?key=" + viewerId);
    let eventSource = new EventSource("https://meme.mmaks.me/realtime/listen?key=" + viewerId);

    eventSource.onmessage = function(event) {
        const parsed = JSON.parse(event.data);
        listener(parsed);
    };
}
