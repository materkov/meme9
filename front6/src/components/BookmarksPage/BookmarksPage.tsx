import {useQuery, useQueryClient} from "@tanstack/react-query";
import React from "react";
import {PostsList} from "../Post/PostsList";
import {getAllFromPosts} from "../../utils/postsList";
import {useGlobals} from "../../store/globals";
import {Link} from "../Link/Link";
import {ApiBookmarks} from "../../api/client";

export function BookmarksPage() {
    const queryClient = useQueryClient();
    const globals = useGlobals();

    const data = useQuery({
        queryKey: ['bookmarks'],
        queryFn: () => (
            ApiBookmarks.List({pageToken: ""}).then(data => {
                for (let item of data.items) {
                    if (item.post) {
                        getAllFromPosts(queryClient, [item.post])
                    }
                }
                return data;
            })
        ),
        enabled: !!globals.viewerId,
    });

    if (!globals.viewerId) {
        return <div><Link href={"/auth"}>Authorize</Link> to see this page</div>
    }

    if (data.status !== "success") {
        return <div>Loading...</div>
    }

    return (
        <div>
            <h1>Bookmarks</h1>
            <PostsList postIds={data.data?.items.map(item => item.post?.id || "")}/>
        </div>
    )
}