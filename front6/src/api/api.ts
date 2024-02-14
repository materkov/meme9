import {cookieAuthToken, getCookie} from "../utils/cookie";

function api<T>(method: string, args: any): Promise<T> {
    let headers: Record<string, string> = {};

    let token = getCookie(cookieAuthToken);
    if (token) {
        headers['authorization'] = 'Bearer ' + token;
    }

    return new Promise((resolve, reject) => {
        // TODO think about this func
        fetch('/api/' + method, {
            credentials: 'omit',
            method: 'POST',
            body: JSON.stringify(args),
            headers: headers,
        })
            .then(r => {
                if (!r.ok) {
                    reject('http error');
                } else if (r.status !== 200) {
                    reject("incorrect http status " + r.status)
                }

                return r.text()
            })
            .then(r => {
                try {
                    const resp = JSON.parse(r);
                    if (resp.error) {
                        reject(resp.error)
                    } else {
                        resolve(resp);
                    }
                } catch (e) {
                    reject('cannot parse json');
                }
            })
            .catch(reject);
    })
}

export class User {
    id = ""
    name = ""
    status = ""
    isFollowing = false
}

export class Post {
    id: string = ""
    userId: string = ""
    date: string = ""
    text: string = ""
    user?: User = undefined

    isLiked: boolean = false
    likesCount: number = 0

    link?: Link = undefined
    poll?: Poll = undefined

    isBookmarked: boolean = false

    photo?: File = undefined
}

export class File {
    url: string = ""
    width: number = 0
    height: number = 0
}

export class Poll {
    id: string = ""
    question = ""
    answers: PollAnswer[] = []
}

export class PollAnswer {
    id: string = ""
    answer: string = ""
    voted: number = 0
    isVoted: boolean = false
}

export class Link {
    title: string = ""
    description: string = ""
    imageUrl: string = ""
    url: string = ""
    domain: string = ""
}

export class PostsList {
    items: Post[] = []
    pageToken: string = ""
}

export class Void {
}

export class PostsAddReq {
    text: string = ""
    pollId: string = ""
    photoId: string = ""
}

export function postsAdd(req: PostsAddReq): Promise<void> {
    return api("meme.api.Posts/Add", req)
}

export enum FeedType {
    UNKNOWN = "",
    DISCOVER = "DISCOVER",
    FEED = "FEED",
}

export class PostsListReq {
    pageToken: string = ""
    count: number = 0
    type: FeedType = FeedType.UNKNOWN
    byUserId: string = ""
    byId: string = ""
}

export function postsList(req: PostsListReq): Promise<PostsList> {
    return api("meme.api.Posts/List", req)
}

export class PostsDeleteReq {
    postId: string = ""
}

export function postsDelete(req: PostsDeleteReq): Promise<Void> {
    return api("meme.api.Posts/Delete", req);
}

export enum LikeAction {
    LIKE = "LIKE",
    UNLIKE = "UNLIKE",
}

export class PostsLikeReq {
    postId: string = ""
    action: LikeAction = LikeAction.LIKE
}

export function postsLike(req: PostsLikeReq): Promise<Void> {
    return api("meme.api.Posts/Like", req);
}

export class UsersListReq {
    userIds: string[] = []
}

export function usersList(req: UsersListReq): Promise<User[]> {
    return api("users.list", req);
}

export class AuthEmailReq {
    email: string = ""
    password: string = ""
}

export class AuthResp {
    token: string = ""
    userId: string = ""
    userName: string = ""
}

export function authLogin(req: AuthEmailReq): Promise<AuthResp> {
    return api("meme.api.Auth/Login", req);
}

export function authRegister(req: AuthEmailReq): Promise<AuthResp> {
    return api("meme.api.Auth/Register", req);
}

export class AuthVk {
    code: string = ""
    redirectUrl: string = ""
}

export function authVK(req: AuthVk): Promise<AuthResp> {
    return api("meme.api.Auth/Vk", req);
}

export class UsersSetStatus {
    status: string = ""
}

export function usersSetStatus(req: UsersSetStatus): Promise<Void> {
    return api("meme.api.Users/SetStatus", req);
}

export enum SubscribeAction {
    NONE = "",
    FOLLOW = "FOLLOW",
    UNFOLLOW = "UNFOLLOW",
}

export class UsersFollowReq {
    targetId: string = ""
    action: SubscribeAction = SubscribeAction.NONE
}

export function usersFollow(req: UsersFollowReq): Promise<Void> {
    return api("meme.api.Users/Follow", req);
}

export class PollsVoteReq {
    pollId: string = ""
    answerIds: string[] = []
}

export function pollsVote(req: PollsVoteReq): Promise<Void> {
    return api("meme.api.Polls/Vote", req);
}

export class PollsDeleteVoteReq {
    pollId: string = ""
}

export function pollsDeleteVote(req: PollsDeleteVoteReq): Promise<Void> {
    return api("meme.api.Polls/DeleteVote", req);
}

export class PollsAddReq {
    question: string = ""
    answers: string[] = []
}

export function pollsAdd(req: PollsAddReq): Promise<Poll> {
    return api("meme.api.Polls/Add", req);
}

export class BookmarksAddReq {
    postId: string = ""
}

export class Bookmark {
    date: string = ""
    post: Post | undefined = undefined
}

export function bookmarksAdd(req: BookmarksAddReq): Promise<Void> {
    return api("meme.api.Bookmarks/Add", req);
}

export function bookmarksRemove(req: BookmarksAddReq): Promise<Void> {
    return api("meme.api.Bookmarks/Remove", req);
}

export class BookmarkListReq {
    pageToken: string = ""
}

export class BookmarkList {
    pageToken: string = ""
    items: Bookmark[] = []
}

export function bookmarksList(req: BookmarkListReq): Promise<BookmarkList> {
    return api("meme.api.Bookmarks/List", req);
}
