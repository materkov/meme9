import {QueryClient} from "@tanstack/react-query";
import * as types from "../api/api";

export function getAllFromPosts(queryClient: QueryClient, posts: types.Post[]) {
    for (let post of posts) {
        queryClient.setQueryData(['post', post.id], post);

        if (post.poll) {
            queryClient.setQueryData(['poll', post.poll.id], post.poll);
        }
        if (post.user) {
            queryClient.setQueryData(['users', post.user.id], post.user);
        }
    }
}
