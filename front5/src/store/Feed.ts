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
            this.posts = [...this.posts, ...data.feed.items];
            this.isLoading = false;
            this.nextCursor = data.feed.nextCursor;
            this.isFirstLoaded = true;
            this.fire();
        });
    }

    public subscribe(callback: () => void) {
        this.callbacks.push(callback)
    }

    public unsubscribe(callback: () => void) {
        this.callbacks = this.callbacks.filter(item => item !== callback)
    }

    public fire() {
        this.version++;
        this.callbacks.forEach(cb => cb());
    }

    public like = (postId: string) => {
        api("/postLike", {id: postId});

        const postsCopy = [...this.posts];
        for (let post of postsCopy) {
            if (post.id == postId) {
                post.isLiked = true;
                post.likesCount = (post.likesCount || 0) + 1;
            }
        }

        this.posts = postsCopy;
        this.fire();
    }

    public unlike = (postId: string) => {
        api("/postUnlike", {id: postId});

        const postsCopy = [...this.posts];
        for (let post of postsCopy) {
            if (post.id == postId) {
                post.isLiked = false;
                post.likesCount = (post.likesCount || 0) - 1;
            }
        }

        this.posts = postsCopy;
        this.fire();
    }

    public reset() {
        this.fetchedCursors = {};
        this.nextCursor = '';
        this.fetch();
    }

    public delete(postId: string) {
        api("/postDelete", {id: postId});

        this.posts = this.posts.filter(post => post.id !== postId);

        this.fire();
    }
}

export const feedStore = new FeedStore();
