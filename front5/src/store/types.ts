export type User = {
    id: string;
    name?: string;
    posts?: Post[];
}

export type Post = {
    id: string;
    userId: string;
    user?: User;
    date: string;
    text: string;
}

//export let apiHost = "http://localhost:8000/api";
export let apiHost = "https://meme.mmaks.me/api";
