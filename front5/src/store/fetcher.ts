export function loadAPI(resources: any): Promise<any> {
    return new Promise((resolve, reject) => {
        const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api" : "http://localhost:8000/api";

        fetch(apiHost, {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
                'authorization': 'Bearer ' + localStorage.getItem('authToken')
            },
            body: JSON.stringify(resources)
        })
            .then(r => r.json())
            .then(r => {
                resolve(r);
            })
    })
}
