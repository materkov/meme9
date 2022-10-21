import {api, Post} from "./types";

export class FeedStore {
    posts: Post[] = [];
    viewerId: string = '';
    nextCursor: string = '';
    isLoading: boolean = false;
    isLoaded: boolean = false;
    isFirstLoaded: boolean = false;

    callbacks: (() => void)[] = [];
    version: number = 0;

    fetchedCursors: { [key: string]: boolean } = {};

    public fetch() {
        if (this.fetchedCursors[this.nextCursor]) return;
        this.fetchedCursors[this.nextCursor] = true;

        this.isLoading = true;
        api("/feed", {cursor: this.nextCursor}).then(data => {
            this.viewerId = data.viewerId;
            this.posts = [...this.posts, ...data.posts];
            this.isLoading = false;
            this.nextCursor = data.nextCursor;
            this.isFirstLoaded = true;
            this.fire();
        });
    }

    public fire() {
        this.version++;
        this.callbacks.forEach(cb => cb());
    }

    public delete(postId: string) {
        api("/postDelete", {id: postId});

        this.posts = this.posts.filter(post => post.id !== postId);

        this.fire();
    }
}

export const feedStore = new FeedStore();
