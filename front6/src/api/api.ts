/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "meme.api";

export enum FeedType {
  FEED = "FEED",
  DISCOVER = "DISCOVER",
  UNRECOGNIZED = "UNRECOGNIZED",
}

export function feedTypeFromJSON(object: any): FeedType {
  switch (object) {
    case 0:
    case "FEED":
      return FeedType.FEED;
    case 1:
    case "DISCOVER":
      return FeedType.DISCOVER;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FeedType.UNRECOGNIZED;
  }
}

export function feedTypeToJSON(object: FeedType): string {
  switch (object) {
    case FeedType.FEED:
      return "FEED";
    case FeedType.DISCOVER:
      return "DISCOVER";
    case FeedType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum PostLikeAction {
  LIKE = "LIKE",
  UNLIKE = "UNLIKE",
  UNRECOGNIZED = "UNRECOGNIZED",
}

export function postLikeActionFromJSON(object: any): PostLikeAction {
  switch (object) {
    case 0:
    case "LIKE":
      return PostLikeAction.LIKE;
    case 1:
    case "UNLIKE":
      return PostLikeAction.UNLIKE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PostLikeAction.UNRECOGNIZED;
  }
}

export function postLikeActionToJSON(object: PostLikeAction): string {
  switch (object) {
    case PostLikeAction.LIKE:
      return "LIKE";
    case PostLikeAction.UNLIKE:
      return "UNLIKE";
    case PostLikeAction.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum SubscribeAction {
  FOLLOW = "FOLLOW",
  UNFOLLOW = "UNFOLLOW",
  UNRECOGNIZED = "UNRECOGNIZED",
}

export function subscribeActionFromJSON(object: any): SubscribeAction {
  switch (object) {
    case 0:
    case "FOLLOW":
      return SubscribeAction.FOLLOW;
    case 1:
    case "UNFOLLOW":
      return SubscribeAction.UNFOLLOW;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SubscribeAction.UNRECOGNIZED;
  }
}

export function subscribeActionToJSON(object: SubscribeAction): string {
  switch (object) {
    case SubscribeAction.FOLLOW:
      return "FOLLOW";
    case SubscribeAction.UNFOLLOW:
      return "UNFOLLOW";
    case SubscribeAction.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Void {
}

export interface AddReq {
  text: string;
  pollId: string;
  photoId: string;
}

export interface Post {
  id: string;
  userId: string;
  date: string;
  text: string;
  user: User | undefined;
  isLiked: boolean;
  likesCount: number;
  link: PostLink | undefined;
  poll: Poll | undefined;
  isBookmarked: boolean;
  isDeleted: boolean;
  photo: File | undefined;
}

export interface File {
  url: string;
  width: number;
  height: number;
}

export interface PostLink {
  url: string;
  title: string;
  description: string;
  imageUrl: string;
  domain: string;
}

export interface PostsList {
  items: Post[];
  pageToken: string;
}

export interface ListReq {
  type: FeedType;
  byUserId: string;
  byId: string;
  count: number;
  pageToken: string;
}

export interface PostsDeleteReq {
  postId: string;
}

export interface PostsLikeReq {
  postId: string;
  action: PostLikeAction;
}

export interface User {
  id: string;
  name: string;
  status: string;
  isFollowing: boolean;
}

export interface UsersList {
  users: User[];
}

export interface UsersListReq {
  userIds: string[];
}

export interface UsersSetStatus {
  status: string;
}

export interface UsersFollowReq {
  targetId: string;
  action: SubscribeAction;
}

export interface CheckAuthReq {
  token: string;
}

export interface EmailReq {
  email: string;
  password: string;
}

export interface AuthResp {
  token: string;
  userId: string;
  userName: string;
}

export interface VkReq {
  code: string;
  redirectUrl: string;
}

export interface PollsList {
  items: Poll[];
}

export interface PollsAddReq {
  question: string;
  answers: string[];
}

export interface Poll {
  id: string;
  question: string;
  answers: PollAnswer[];
}

export interface PollAnswer {
  id: string;
  answer: string;
  votedCount: number;
  isVoted: boolean;
}

export interface PollsListReq {
  ids: string[];
}

export interface PollsVoteReq {
  pollId: string;
  answerIds: string[];
}

export interface PollsDeleteVoteReq {
  pollId: string;
}

export interface BookmarksAddReq {
  postId: string;
}

export interface Bookmark {
  date: string;
  post: Post | undefined;
}

export interface BookmarkListReq {
  pageToken: string;
}

export interface BookmarkList {
  items: Bookmark[];
  pageToken: string;
}

export interface UploadReq {
  photoBytes: Uint8Array;
}

export interface UploadResp {
  uploadToken: string;
}

export interface ResizeReq {
  imageUrl: string;
}

export interface ResizeResp {
  image: Uint8Array;
}

export interface GetEventsReq {
  userId: number;
  timeoutMs: number;
}

export interface GetEventsResp {
  data: string;
}

export interface SendEventReq {
  userId: number;
  payload: string;
}

function createBaseVoid(): Void {
  return {};
}

export const Void = {
  fromJSON(_: any): Void {
    return {};
  },

  toJSON(_: Void): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<Void>, I>>(base?: I): Void {
    return Void.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Void>, I>>(_: I): Void {
    const message = createBaseVoid();
    return message;
  },
};

function createBaseAddReq(): AddReq {
  return { text: "", pollId: "", photoId: "" };
}

export const AddReq = {
  fromJSON(object: any): AddReq {
    return {
      text: isSet(object.text) ? globalThis.String(object.text) : "",
      pollId: isSet(object.pollId) ? globalThis.String(object.pollId) : "",
      photoId: isSet(object.photoId) ? globalThis.String(object.photoId) : "",
    };
  },

  toJSON(message: AddReq): unknown {
    const obj: any = {};
    if (message.text !== "") {
      obj.text = message.text;
    }
    if (message.pollId !== "") {
      obj.pollId = message.pollId;
    }
    if (message.photoId !== "") {
      obj.photoId = message.photoId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AddReq>, I>>(base?: I): AddReq {
    return AddReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AddReq>, I>>(object: I): AddReq {
    const message = createBaseAddReq();
    message.text = object.text ?? "";
    message.pollId = object.pollId ?? "";
    message.photoId = object.photoId ?? "";
    return message;
  },
};

function createBasePost(): Post {
  return {
    id: "",
    userId: "",
    date: "",
    text: "",
    user: undefined,
    isLiked: false,
    likesCount: 0,
    link: undefined,
    poll: undefined,
    isBookmarked: false,
    isDeleted: false,
    photo: undefined,
  };
}

export const Post = {
  fromJSON(object: any): Post {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      userId: isSet(object.userId) ? globalThis.String(object.userId) : "",
      date: isSet(object.date) ? globalThis.String(object.date) : "",
      text: isSet(object.text) ? globalThis.String(object.text) : "",
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
      isLiked: isSet(object.isLiked) ? globalThis.Boolean(object.isLiked) : false,
      likesCount: isSet(object.likesCount) ? globalThis.Number(object.likesCount) : 0,
      link: isSet(object.link) ? PostLink.fromJSON(object.link) : undefined,
      poll: isSet(object.poll) ? Poll.fromJSON(object.poll) : undefined,
      isBookmarked: isSet(object.isBookmarked) ? globalThis.Boolean(object.isBookmarked) : false,
      isDeleted: isSet(object.isDeleted) ? globalThis.Boolean(object.isDeleted) : false,
      photo: isSet(object.photo) ? File.fromJSON(object.photo) : undefined,
    };
  },

  toJSON(message: Post): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    if (message.date !== "") {
      obj.date = message.date;
    }
    if (message.text !== "") {
      obj.text = message.text;
    }
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    if (message.isLiked === true) {
      obj.isLiked = message.isLiked;
    }
    if (message.likesCount !== 0) {
      obj.likesCount = Math.round(message.likesCount);
    }
    if (message.link !== undefined) {
      obj.link = PostLink.toJSON(message.link);
    }
    if (message.poll !== undefined) {
      obj.poll = Poll.toJSON(message.poll);
    }
    if (message.isBookmarked === true) {
      obj.isBookmarked = message.isBookmarked;
    }
    if (message.isDeleted === true) {
      obj.isDeleted = message.isDeleted;
    }
    if (message.photo !== undefined) {
      obj.photo = File.toJSON(message.photo);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Post>, I>>(base?: I): Post {
    return Post.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Post>, I>>(object: I): Post {
    const message = createBasePost();
    message.id = object.id ?? "";
    message.userId = object.userId ?? "";
    message.date = object.date ?? "";
    message.text = object.text ?? "";
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    message.isLiked = object.isLiked ?? false;
    message.likesCount = object.likesCount ?? 0;
    message.link = (object.link !== undefined && object.link !== null) ? PostLink.fromPartial(object.link) : undefined;
    message.poll = (object.poll !== undefined && object.poll !== null) ? Poll.fromPartial(object.poll) : undefined;
    message.isBookmarked = object.isBookmarked ?? false;
    message.isDeleted = object.isDeleted ?? false;
    message.photo = (object.photo !== undefined && object.photo !== null) ? File.fromPartial(object.photo) : undefined;
    return message;
  },
};

function createBaseFile(): File {
  return { url: "", width: 0, height: 0 };
}

export const File = {
  fromJSON(object: any): File {
    return {
      url: isSet(object.url) ? globalThis.String(object.url) : "",
      width: isSet(object.width) ? globalThis.Number(object.width) : 0,
      height: isSet(object.height) ? globalThis.Number(object.height) : 0,
    };
  },

  toJSON(message: File): unknown {
    const obj: any = {};
    if (message.url !== "") {
      obj.url = message.url;
    }
    if (message.width !== 0) {
      obj.width = Math.round(message.width);
    }
    if (message.height !== 0) {
      obj.height = Math.round(message.height);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<File>, I>>(base?: I): File {
    return File.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<File>, I>>(object: I): File {
    const message = createBaseFile();
    message.url = object.url ?? "";
    message.width = object.width ?? 0;
    message.height = object.height ?? 0;
    return message;
  },
};

function createBasePostLink(): PostLink {
  return { url: "", title: "", description: "", imageUrl: "", domain: "" };
}

export const PostLink = {
  fromJSON(object: any): PostLink {
    return {
      url: isSet(object.url) ? globalThis.String(object.url) : "",
      title: isSet(object.title) ? globalThis.String(object.title) : "",
      description: isSet(object.description) ? globalThis.String(object.description) : "",
      imageUrl: isSet(object.imageUrl) ? globalThis.String(object.imageUrl) : "",
      domain: isSet(object.domain) ? globalThis.String(object.domain) : "",
    };
  },

  toJSON(message: PostLink): unknown {
    const obj: any = {};
    if (message.url !== "") {
      obj.url = message.url;
    }
    if (message.title !== "") {
      obj.title = message.title;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.imageUrl !== "") {
      obj.imageUrl = message.imageUrl;
    }
    if (message.domain !== "") {
      obj.domain = message.domain;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PostLink>, I>>(base?: I): PostLink {
    return PostLink.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PostLink>, I>>(object: I): PostLink {
    const message = createBasePostLink();
    message.url = object.url ?? "";
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.imageUrl = object.imageUrl ?? "";
    message.domain = object.domain ?? "";
    return message;
  },
};

function createBasePostsList(): PostsList {
  return { items: [], pageToken: "" };
}

export const PostsList = {
  fromJSON(object: any): PostsList {
    return {
      items: globalThis.Array.isArray(object?.items) ? object.items.map((e: any) => Post.fromJSON(e)) : [],
      pageToken: isSet(object.pageToken) ? globalThis.String(object.pageToken) : "",
    };
  },

  toJSON(message: PostsList): unknown {
    const obj: any = {};
    if (message.items?.length) {
      obj.items = message.items.map((e) => Post.toJSON(e));
    }
    if (message.pageToken !== "") {
      obj.pageToken = message.pageToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PostsList>, I>>(base?: I): PostsList {
    return PostsList.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PostsList>, I>>(object: I): PostsList {
    const message = createBasePostsList();
    message.items = object.items?.map((e) => Post.fromPartial(e)) || [];
    message.pageToken = object.pageToken ?? "";
    return message;
  },
};

function createBaseListReq(): ListReq {
  return { type: FeedType.FEED, byUserId: "", byId: "", count: 0, pageToken: "" };
}

export const ListReq = {
  fromJSON(object: any): ListReq {
    return {
      type: isSet(object.type) ? feedTypeFromJSON(object.type) : FeedType.FEED,
      byUserId: isSet(object.byUserId) ? globalThis.String(object.byUserId) : "",
      byId: isSet(object.byId) ? globalThis.String(object.byId) : "",
      count: isSet(object.count) ? globalThis.Number(object.count) : 0,
      pageToken: isSet(object.pageToken) ? globalThis.String(object.pageToken) : "",
    };
  },

  toJSON(message: ListReq): unknown {
    const obj: any = {};
    if (message.type !== FeedType.FEED) {
      obj.type = feedTypeToJSON(message.type);
    }
    if (message.byUserId !== "") {
      obj.byUserId = message.byUserId;
    }
    if (message.byId !== "") {
      obj.byId = message.byId;
    }
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    if (message.pageToken !== "") {
      obj.pageToken = message.pageToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListReq>, I>>(base?: I): ListReq {
    return ListReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListReq>, I>>(object: I): ListReq {
    const message = createBaseListReq();
    message.type = object.type ?? FeedType.FEED;
    message.byUserId = object.byUserId ?? "";
    message.byId = object.byId ?? "";
    message.count = object.count ?? 0;
    message.pageToken = object.pageToken ?? "";
    return message;
  },
};

function createBasePostsDeleteReq(): PostsDeleteReq {
  return { postId: "" };
}

export const PostsDeleteReq = {
  fromJSON(object: any): PostsDeleteReq {
    return { postId: isSet(object.postId) ? globalThis.String(object.postId) : "" };
  },

  toJSON(message: PostsDeleteReq): unknown {
    const obj: any = {};
    if (message.postId !== "") {
      obj.postId = message.postId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PostsDeleteReq>, I>>(base?: I): PostsDeleteReq {
    return PostsDeleteReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PostsDeleteReq>, I>>(object: I): PostsDeleteReq {
    const message = createBasePostsDeleteReq();
    message.postId = object.postId ?? "";
    return message;
  },
};

function createBasePostsLikeReq(): PostsLikeReq {
  return { postId: "", action: PostLikeAction.LIKE };
}

export const PostsLikeReq = {
  fromJSON(object: any): PostsLikeReq {
    return {
      postId: isSet(object.postId) ? globalThis.String(object.postId) : "",
      action: isSet(object.action) ? postLikeActionFromJSON(object.action) : PostLikeAction.LIKE,
    };
  },

  toJSON(message: PostsLikeReq): unknown {
    const obj: any = {};
    if (message.postId !== "") {
      obj.postId = message.postId;
    }
    if (message.action !== PostLikeAction.LIKE) {
      obj.action = postLikeActionToJSON(message.action);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PostsLikeReq>, I>>(base?: I): PostsLikeReq {
    return PostsLikeReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PostsLikeReq>, I>>(object: I): PostsLikeReq {
    const message = createBasePostsLikeReq();
    message.postId = object.postId ?? "";
    message.action = object.action ?? PostLikeAction.LIKE;
    return message;
  },
};

function createBaseUser(): User {
  return { id: "", name: "", status: "", isFollowing: false };
}

export const User = {
  fromJSON(object: any): User {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      status: isSet(object.status) ? globalThis.String(object.status) : "",
      isFollowing: isSet(object.isFollowing) ? globalThis.Boolean(object.isFollowing) : false,
    };
  },

  toJSON(message: User): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.status !== "") {
      obj.status = message.status;
    }
    if (message.isFollowing === true) {
      obj.isFollowing = message.isFollowing;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<User>, I>>(base?: I): User {
    return User.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<User>, I>>(object: I): User {
    const message = createBaseUser();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.status = object.status ?? "";
    message.isFollowing = object.isFollowing ?? false;
    return message;
  },
};

function createBaseUsersList(): UsersList {
  return { users: [] };
}

export const UsersList = {
  fromJSON(object: any): UsersList {
    return { users: globalThis.Array.isArray(object?.users) ? object.users.map((e: any) => User.fromJSON(e)) : [] };
  },

  toJSON(message: UsersList): unknown {
    const obj: any = {};
    if (message.users?.length) {
      obj.users = message.users.map((e) => User.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UsersList>, I>>(base?: I): UsersList {
    return UsersList.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UsersList>, I>>(object: I): UsersList {
    const message = createBaseUsersList();
    message.users = object.users?.map((e) => User.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUsersListReq(): UsersListReq {
  return { userIds: [] };
}

export const UsersListReq = {
  fromJSON(object: any): UsersListReq {
    return {
      userIds: globalThis.Array.isArray(object?.userIds) ? object.userIds.map((e: any) => globalThis.String(e)) : [],
    };
  },

  toJSON(message: UsersListReq): unknown {
    const obj: any = {};
    if (message.userIds?.length) {
      obj.userIds = message.userIds;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UsersListReq>, I>>(base?: I): UsersListReq {
    return UsersListReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UsersListReq>, I>>(object: I): UsersListReq {
    const message = createBaseUsersListReq();
    message.userIds = object.userIds?.map((e) => e) || [];
    return message;
  },
};

function createBaseUsersSetStatus(): UsersSetStatus {
  return { status: "" };
}

export const UsersSetStatus = {
  fromJSON(object: any): UsersSetStatus {
    return { status: isSet(object.status) ? globalThis.String(object.status) : "" };
  },

  toJSON(message: UsersSetStatus): unknown {
    const obj: any = {};
    if (message.status !== "") {
      obj.status = message.status;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UsersSetStatus>, I>>(base?: I): UsersSetStatus {
    return UsersSetStatus.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UsersSetStatus>, I>>(object: I): UsersSetStatus {
    const message = createBaseUsersSetStatus();
    message.status = object.status ?? "";
    return message;
  },
};

function createBaseUsersFollowReq(): UsersFollowReq {
  return { targetId: "", action: SubscribeAction.FOLLOW };
}

export const UsersFollowReq = {
  fromJSON(object: any): UsersFollowReq {
    return {
      targetId: isSet(object.targetId) ? globalThis.String(object.targetId) : "",
      action: isSet(object.action) ? subscribeActionFromJSON(object.action) : SubscribeAction.FOLLOW,
    };
  },

  toJSON(message: UsersFollowReq): unknown {
    const obj: any = {};
    if (message.targetId !== "") {
      obj.targetId = message.targetId;
    }
    if (message.action !== SubscribeAction.FOLLOW) {
      obj.action = subscribeActionToJSON(message.action);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UsersFollowReq>, I>>(base?: I): UsersFollowReq {
    return UsersFollowReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UsersFollowReq>, I>>(object: I): UsersFollowReq {
    const message = createBaseUsersFollowReq();
    message.targetId = object.targetId ?? "";
    message.action = object.action ?? SubscribeAction.FOLLOW;
    return message;
  },
};

function createBaseCheckAuthReq(): CheckAuthReq {
  return { token: "" };
}

export const CheckAuthReq = {
  fromJSON(object: any): CheckAuthReq {
    return { token: isSet(object.token) ? globalThis.String(object.token) : "" };
  },

  toJSON(message: CheckAuthReq): unknown {
    const obj: any = {};
    if (message.token !== "") {
      obj.token = message.token;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAuthReq>, I>>(base?: I): CheckAuthReq {
    return CheckAuthReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAuthReq>, I>>(object: I): CheckAuthReq {
    const message = createBaseCheckAuthReq();
    message.token = object.token ?? "";
    return message;
  },
};

function createBaseEmailReq(): EmailReq {
  return { email: "", password: "" };
}

export const EmailReq = {
  fromJSON(object: any): EmailReq {
    return {
      email: isSet(object.email) ? globalThis.String(object.email) : "",
      password: isSet(object.password) ? globalThis.String(object.password) : "",
    };
  },

  toJSON(message: EmailReq): unknown {
    const obj: any = {};
    if (message.email !== "") {
      obj.email = message.email;
    }
    if (message.password !== "") {
      obj.password = message.password;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EmailReq>, I>>(base?: I): EmailReq {
    return EmailReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EmailReq>, I>>(object: I): EmailReq {
    const message = createBaseEmailReq();
    message.email = object.email ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseAuthResp(): AuthResp {
  return { token: "", userId: "", userName: "" };
}

export const AuthResp = {
  fromJSON(object: any): AuthResp {
    return {
      token: isSet(object.token) ? globalThis.String(object.token) : "",
      userId: isSet(object.userId) ? globalThis.String(object.userId) : "",
      userName: isSet(object.userName) ? globalThis.String(object.userName) : "",
    };
  },

  toJSON(message: AuthResp): unknown {
    const obj: any = {};
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    if (message.userName !== "") {
      obj.userName = message.userName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AuthResp>, I>>(base?: I): AuthResp {
    return AuthResp.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AuthResp>, I>>(object: I): AuthResp {
    const message = createBaseAuthResp();
    message.token = object.token ?? "";
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    return message;
  },
};

function createBaseVkReq(): VkReq {
  return { code: "", redirectUrl: "" };
}

export const VkReq = {
  fromJSON(object: any): VkReq {
    return {
      code: isSet(object.code) ? globalThis.String(object.code) : "",
      redirectUrl: isSet(object.redirectUrl) ? globalThis.String(object.redirectUrl) : "",
    };
  },

  toJSON(message: VkReq): unknown {
    const obj: any = {};
    if (message.code !== "") {
      obj.code = message.code;
    }
    if (message.redirectUrl !== "") {
      obj.redirectUrl = message.redirectUrl;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<VkReq>, I>>(base?: I): VkReq {
    return VkReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<VkReq>, I>>(object: I): VkReq {
    const message = createBaseVkReq();
    message.code = object.code ?? "";
    message.redirectUrl = object.redirectUrl ?? "";
    return message;
  },
};

function createBasePollsList(): PollsList {
  return { items: [] };
}

export const PollsList = {
  fromJSON(object: any): PollsList {
    return { items: globalThis.Array.isArray(object?.items) ? object.items.map((e: any) => Poll.fromJSON(e)) : [] };
  },

  toJSON(message: PollsList): unknown {
    const obj: any = {};
    if (message.items?.length) {
      obj.items = message.items.map((e) => Poll.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PollsList>, I>>(base?: I): PollsList {
    return PollsList.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PollsList>, I>>(object: I): PollsList {
    const message = createBasePollsList();
    message.items = object.items?.map((e) => Poll.fromPartial(e)) || [];
    return message;
  },
};

function createBasePollsAddReq(): PollsAddReq {
  return { question: "", answers: [] };
}

export const PollsAddReq = {
  fromJSON(object: any): PollsAddReq {
    return {
      question: isSet(object.question) ? globalThis.String(object.question) : "",
      answers: globalThis.Array.isArray(object?.answers) ? object.answers.map((e: any) => globalThis.String(e)) : [],
    };
  },

  toJSON(message: PollsAddReq): unknown {
    const obj: any = {};
    if (message.question !== "") {
      obj.question = message.question;
    }
    if (message.answers?.length) {
      obj.answers = message.answers;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PollsAddReq>, I>>(base?: I): PollsAddReq {
    return PollsAddReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PollsAddReq>, I>>(object: I): PollsAddReq {
    const message = createBasePollsAddReq();
    message.question = object.question ?? "";
    message.answers = object.answers?.map((e) => e) || [];
    return message;
  },
};

function createBasePoll(): Poll {
  return { id: "", question: "", answers: [] };
}

export const Poll = {
  fromJSON(object: any): Poll {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      question: isSet(object.question) ? globalThis.String(object.question) : "",
      answers: globalThis.Array.isArray(object?.answers) ? object.answers.map((e: any) => PollAnswer.fromJSON(e)) : [],
    };
  },

  toJSON(message: Poll): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.question !== "") {
      obj.question = message.question;
    }
    if (message.answers?.length) {
      obj.answers = message.answers.map((e) => PollAnswer.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Poll>, I>>(base?: I): Poll {
    return Poll.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Poll>, I>>(object: I): Poll {
    const message = createBasePoll();
    message.id = object.id ?? "";
    message.question = object.question ?? "";
    message.answers = object.answers?.map((e) => PollAnswer.fromPartial(e)) || [];
    return message;
  },
};

function createBasePollAnswer(): PollAnswer {
  return { id: "", answer: "", votedCount: 0, isVoted: false };
}

export const PollAnswer = {
  fromJSON(object: any): PollAnswer {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      answer: isSet(object.answer) ? globalThis.String(object.answer) : "",
      votedCount: isSet(object.votedCount) ? globalThis.Number(object.votedCount) : 0,
      isVoted: isSet(object.isVoted) ? globalThis.Boolean(object.isVoted) : false,
    };
  },

  toJSON(message: PollAnswer): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.answer !== "") {
      obj.answer = message.answer;
    }
    if (message.votedCount !== 0) {
      obj.votedCount = Math.round(message.votedCount);
    }
    if (message.isVoted === true) {
      obj.isVoted = message.isVoted;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PollAnswer>, I>>(base?: I): PollAnswer {
    return PollAnswer.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PollAnswer>, I>>(object: I): PollAnswer {
    const message = createBasePollAnswer();
    message.id = object.id ?? "";
    message.answer = object.answer ?? "";
    message.votedCount = object.votedCount ?? 0;
    message.isVoted = object.isVoted ?? false;
    return message;
  },
};

function createBasePollsListReq(): PollsListReq {
  return { ids: [] };
}

export const PollsListReq = {
  fromJSON(object: any): PollsListReq {
    return { ids: globalThis.Array.isArray(object?.ids) ? object.ids.map((e: any) => globalThis.String(e)) : [] };
  },

  toJSON(message: PollsListReq): unknown {
    const obj: any = {};
    if (message.ids?.length) {
      obj.ids = message.ids;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PollsListReq>, I>>(base?: I): PollsListReq {
    return PollsListReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PollsListReq>, I>>(object: I): PollsListReq {
    const message = createBasePollsListReq();
    message.ids = object.ids?.map((e) => e) || [];
    return message;
  },
};

function createBasePollsVoteReq(): PollsVoteReq {
  return { pollId: "", answerIds: [] };
}

export const PollsVoteReq = {
  fromJSON(object: any): PollsVoteReq {
    return {
      pollId: isSet(object.pollId) ? globalThis.String(object.pollId) : "",
      answerIds: globalThis.Array.isArray(object?.answerIds)
        ? object.answerIds.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: PollsVoteReq): unknown {
    const obj: any = {};
    if (message.pollId !== "") {
      obj.pollId = message.pollId;
    }
    if (message.answerIds?.length) {
      obj.answerIds = message.answerIds;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PollsVoteReq>, I>>(base?: I): PollsVoteReq {
    return PollsVoteReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PollsVoteReq>, I>>(object: I): PollsVoteReq {
    const message = createBasePollsVoteReq();
    message.pollId = object.pollId ?? "";
    message.answerIds = object.answerIds?.map((e) => e) || [];
    return message;
  },
};

function createBasePollsDeleteVoteReq(): PollsDeleteVoteReq {
  return { pollId: "" };
}

export const PollsDeleteVoteReq = {
  fromJSON(object: any): PollsDeleteVoteReq {
    return { pollId: isSet(object.pollId) ? globalThis.String(object.pollId) : "" };
  },

  toJSON(message: PollsDeleteVoteReq): unknown {
    const obj: any = {};
    if (message.pollId !== "") {
      obj.pollId = message.pollId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PollsDeleteVoteReq>, I>>(base?: I): PollsDeleteVoteReq {
    return PollsDeleteVoteReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PollsDeleteVoteReq>, I>>(object: I): PollsDeleteVoteReq {
    const message = createBasePollsDeleteVoteReq();
    message.pollId = object.pollId ?? "";
    return message;
  },
};

function createBaseBookmarksAddReq(): BookmarksAddReq {
  return { postId: "" };
}

export const BookmarksAddReq = {
  fromJSON(object: any): BookmarksAddReq {
    return { postId: isSet(object.postId) ? globalThis.String(object.postId) : "" };
  },

  toJSON(message: BookmarksAddReq): unknown {
    const obj: any = {};
    if (message.postId !== "") {
      obj.postId = message.postId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<BookmarksAddReq>, I>>(base?: I): BookmarksAddReq {
    return BookmarksAddReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<BookmarksAddReq>, I>>(object: I): BookmarksAddReq {
    const message = createBaseBookmarksAddReq();
    message.postId = object.postId ?? "";
    return message;
  },
};

function createBaseBookmark(): Bookmark {
  return { date: "", post: undefined };
}

export const Bookmark = {
  fromJSON(object: any): Bookmark {
    return {
      date: isSet(object.date) ? globalThis.String(object.date) : "",
      post: isSet(object.post) ? Post.fromJSON(object.post) : undefined,
    };
  },

  toJSON(message: Bookmark): unknown {
    const obj: any = {};
    if (message.date !== "") {
      obj.date = message.date;
    }
    if (message.post !== undefined) {
      obj.post = Post.toJSON(message.post);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Bookmark>, I>>(base?: I): Bookmark {
    return Bookmark.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Bookmark>, I>>(object: I): Bookmark {
    const message = createBaseBookmark();
    message.date = object.date ?? "";
    message.post = (object.post !== undefined && object.post !== null) ? Post.fromPartial(object.post) : undefined;
    return message;
  },
};

function createBaseBookmarkListReq(): BookmarkListReq {
  return { pageToken: "" };
}

export const BookmarkListReq = {
  fromJSON(object: any): BookmarkListReq {
    return { pageToken: isSet(object.pageToken) ? globalThis.String(object.pageToken) : "" };
  },

  toJSON(message: BookmarkListReq): unknown {
    const obj: any = {};
    if (message.pageToken !== "") {
      obj.pageToken = message.pageToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<BookmarkListReq>, I>>(base?: I): BookmarkListReq {
    return BookmarkListReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<BookmarkListReq>, I>>(object: I): BookmarkListReq {
    const message = createBaseBookmarkListReq();
    message.pageToken = object.pageToken ?? "";
    return message;
  },
};

function createBaseBookmarkList(): BookmarkList {
  return { items: [], pageToken: "" };
}

export const BookmarkList = {
  fromJSON(object: any): BookmarkList {
    return {
      items: globalThis.Array.isArray(object?.items) ? object.items.map((e: any) => Bookmark.fromJSON(e)) : [],
      pageToken: isSet(object.pageToken) ? globalThis.String(object.pageToken) : "",
    };
  },

  toJSON(message: BookmarkList): unknown {
    const obj: any = {};
    if (message.items?.length) {
      obj.items = message.items.map((e) => Bookmark.toJSON(e));
    }
    if (message.pageToken !== "") {
      obj.pageToken = message.pageToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<BookmarkList>, I>>(base?: I): BookmarkList {
    return BookmarkList.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<BookmarkList>, I>>(object: I): BookmarkList {
    const message = createBaseBookmarkList();
    message.items = object.items?.map((e) => Bookmark.fromPartial(e)) || [];
    message.pageToken = object.pageToken ?? "";
    return message;
  },
};

function createBaseUploadReq(): UploadReq {
  return { photoBytes: new Uint8Array(0) };
}

export const UploadReq = {
  fromJSON(object: any): UploadReq {
    return { photoBytes: isSet(object.photoBytes) ? bytesFromBase64(object.photoBytes) : new Uint8Array(0) };
  },

  toJSON(message: UploadReq): unknown {
    const obj: any = {};
    if (message.photoBytes.length !== 0) {
      obj.photoBytes = base64FromBytes(message.photoBytes);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UploadReq>, I>>(base?: I): UploadReq {
    return UploadReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UploadReq>, I>>(object: I): UploadReq {
    const message = createBaseUploadReq();
    message.photoBytes = object.photoBytes ?? new Uint8Array(0);
    return message;
  },
};

function createBaseUploadResp(): UploadResp {
  return { uploadToken: "" };
}

export const UploadResp = {
  fromJSON(object: any): UploadResp {
    return { uploadToken: isSet(object.uploadToken) ? globalThis.String(object.uploadToken) : "" };
  },

  toJSON(message: UploadResp): unknown {
    const obj: any = {};
    if (message.uploadToken !== "") {
      obj.uploadToken = message.uploadToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UploadResp>, I>>(base?: I): UploadResp {
    return UploadResp.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UploadResp>, I>>(object: I): UploadResp {
    const message = createBaseUploadResp();
    message.uploadToken = object.uploadToken ?? "";
    return message;
  },
};

function createBaseResizeReq(): ResizeReq {
  return { imageUrl: "" };
}

export const ResizeReq = {
  fromJSON(object: any): ResizeReq {
    return { imageUrl: isSet(object.imageUrl) ? globalThis.String(object.imageUrl) : "" };
  },

  toJSON(message: ResizeReq): unknown {
    const obj: any = {};
    if (message.imageUrl !== "") {
      obj.imageUrl = message.imageUrl;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ResizeReq>, I>>(base?: I): ResizeReq {
    return ResizeReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ResizeReq>, I>>(object: I): ResizeReq {
    const message = createBaseResizeReq();
    message.imageUrl = object.imageUrl ?? "";
    return message;
  },
};

function createBaseResizeResp(): ResizeResp {
  return { image: new Uint8Array(0) };
}

export const ResizeResp = {
  fromJSON(object: any): ResizeResp {
    return { image: isSet(object.image) ? bytesFromBase64(object.image) : new Uint8Array(0) };
  },

  toJSON(message: ResizeResp): unknown {
    const obj: any = {};
    if (message.image.length !== 0) {
      obj.image = base64FromBytes(message.image);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ResizeResp>, I>>(base?: I): ResizeResp {
    return ResizeResp.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ResizeResp>, I>>(object: I): ResizeResp {
    const message = createBaseResizeResp();
    message.image = object.image ?? new Uint8Array(0);
    return message;
  },
};

function createBaseGetEventsReq(): GetEventsReq {
  return { userId: 0, timeoutMs: 0 };
}

export const GetEventsReq = {
  fromJSON(object: any): GetEventsReq {
    return {
      userId: isSet(object.userId) ? globalThis.Number(object.userId) : 0,
      timeoutMs: isSet(object.timeoutMs) ? globalThis.Number(object.timeoutMs) : 0,
    };
  },

  toJSON(message: GetEventsReq): unknown {
    const obj: any = {};
    if (message.userId !== 0) {
      obj.userId = Math.round(message.userId);
    }
    if (message.timeoutMs !== 0) {
      obj.timeoutMs = Math.round(message.timeoutMs);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetEventsReq>, I>>(base?: I): GetEventsReq {
    return GetEventsReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetEventsReq>, I>>(object: I): GetEventsReq {
    const message = createBaseGetEventsReq();
    message.userId = object.userId ?? 0;
    message.timeoutMs = object.timeoutMs ?? 0;
    return message;
  },
};

function createBaseGetEventsResp(): GetEventsResp {
  return { data: "" };
}

export const GetEventsResp = {
  fromJSON(object: any): GetEventsResp {
    return { data: isSet(object.data) ? globalThis.String(object.data) : "" };
  },

  toJSON(message: GetEventsResp): unknown {
    const obj: any = {};
    if (message.data !== "") {
      obj.data = message.data;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetEventsResp>, I>>(base?: I): GetEventsResp {
    return GetEventsResp.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetEventsResp>, I>>(object: I): GetEventsResp {
    const message = createBaseGetEventsResp();
    message.data = object.data ?? "";
    return message;
  },
};

function createBaseSendEventReq(): SendEventReq {
  return { userId: 0, payload: "" };
}

export const SendEventReq = {
  fromJSON(object: any): SendEventReq {
    return {
      userId: isSet(object.userId) ? globalThis.Number(object.userId) : 0,
      payload: isSet(object.payload) ? globalThis.String(object.payload) : "",
    };
  },

  toJSON(message: SendEventReq): unknown {
    const obj: any = {};
    if (message.userId !== 0) {
      obj.userId = Math.round(message.userId);
    }
    if (message.payload !== "") {
      obj.payload = message.payload;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SendEventReq>, I>>(base?: I): SendEventReq {
    return SendEventReq.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SendEventReq>, I>>(object: I): SendEventReq {
    const message = createBaseSendEventReq();
    message.userId = object.userId ?? 0;
    message.payload = object.payload ?? "";
    return message;
  },
};

export interface Posts {
  Add(request: AddReq): Promise<Post>;
  List(request: ListReq): Promise<PostsList>;
  Delete(request: PostsDeleteReq): Promise<Void>;
  Like(request: PostsLikeReq): Promise<Void>;
}

export const PostsServiceName = "meme.api.Posts";
export class PostsClientImpl implements Posts {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || PostsServiceName;
    this.rpc = rpc;
    this.Add = this.Add.bind(this);
    this.List = this.List.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Like = this.Like.bind(this);
  }
  Add(request: AddReq): Promise<Post> {
    const data = AddReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Add", data);
    return promise.then((data) => Post.fromJSON(data));
  }

  List(request: ListReq): Promise<PostsList> {
    const data = ListReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "List", data);
    return promise.then((data) => PostsList.fromJSON(data));
  }

  Delete(request: PostsDeleteReq): Promise<Void> {
    const data = PostsDeleteReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => Void.fromJSON(data));
  }

  Like(request: PostsLikeReq): Promise<Void> {
    const data = PostsLikeReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Like", data);
    return promise.then((data) => Void.fromJSON(data));
  }
}

export interface Users {
  List(request: UsersListReq): Promise<UsersList>;
  SetStatus(request: UsersSetStatus): Promise<Void>;
  Follow(request: UsersFollowReq): Promise<Void>;
}

export const UsersServiceName = "meme.api.Users";
export class UsersClientImpl implements Users {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || UsersServiceName;
    this.rpc = rpc;
    this.List = this.List.bind(this);
    this.SetStatus = this.SetStatus.bind(this);
    this.Follow = this.Follow.bind(this);
  }
  List(request: UsersListReq): Promise<UsersList> {
    const data = UsersListReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "List", data);
    return promise.then((data) => UsersList.fromJSON(data));
  }

  SetStatus(request: UsersSetStatus): Promise<Void> {
    const data = UsersSetStatus.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "SetStatus", data);
    return promise.then((data) => Void.fromJSON(data));
  }

  Follow(request: UsersFollowReq): Promise<Void> {
    const data = UsersFollowReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Follow", data);
    return promise.then((data) => Void.fromJSON(data));
  }
}

export interface Auth {
  Login(request: EmailReq): Promise<AuthResp>;
  Register(request: EmailReq): Promise<AuthResp>;
  Vk(request: VkReq): Promise<AuthResp>;
  CheckAuth(request: CheckAuthReq): Promise<AuthResp>;
}

export const AuthServiceName = "meme.api.Auth";
export class AuthClientImpl implements Auth {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || AuthServiceName;
    this.rpc = rpc;
    this.Login = this.Login.bind(this);
    this.Register = this.Register.bind(this);
    this.Vk = this.Vk.bind(this);
    this.CheckAuth = this.CheckAuth.bind(this);
  }
  Login(request: EmailReq): Promise<AuthResp> {
    const data = EmailReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Login", data);
    return promise.then((data) => AuthResp.fromJSON(data));
  }

  Register(request: EmailReq): Promise<AuthResp> {
    const data = EmailReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Register", data);
    return promise.then((data) => AuthResp.fromJSON(data));
  }

  Vk(request: VkReq): Promise<AuthResp> {
    const data = VkReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Vk", data);
    return promise.then((data) => AuthResp.fromJSON(data));
  }

  CheckAuth(request: CheckAuthReq): Promise<AuthResp> {
    const data = CheckAuthReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "CheckAuth", data);
    return promise.then((data) => AuthResp.fromJSON(data));
  }
}

export interface Polls {
  Add(request: PollsAddReq): Promise<Poll>;
  List(request: PollsListReq): Promise<PollsList>;
  Vote(request: PollsVoteReq): Promise<Void>;
  DeleteVote(request: PollsDeleteVoteReq): Promise<Void>;
}

export const PollsServiceName = "meme.api.Polls";
export class PollsClientImpl implements Polls {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || PollsServiceName;
    this.rpc = rpc;
    this.Add = this.Add.bind(this);
    this.List = this.List.bind(this);
    this.Vote = this.Vote.bind(this);
    this.DeleteVote = this.DeleteVote.bind(this);
  }
  Add(request: PollsAddReq): Promise<Poll> {
    const data = PollsAddReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Add", data);
    return promise.then((data) => Poll.fromJSON(data));
  }

  List(request: PollsListReq): Promise<PollsList> {
    const data = PollsListReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "List", data);
    return promise.then((data) => PollsList.fromJSON(data));
  }

  Vote(request: PollsVoteReq): Promise<Void> {
    const data = PollsVoteReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Vote", data);
    return promise.then((data) => Void.fromJSON(data));
  }

  DeleteVote(request: PollsDeleteVoteReq): Promise<Void> {
    const data = PollsDeleteVoteReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "DeleteVote", data);
    return promise.then((data) => Void.fromJSON(data));
  }
}

export interface Bookmarks {
  Add(request: BookmarksAddReq): Promise<Void>;
  Remove(request: BookmarksAddReq): Promise<Void>;
  List(request: BookmarkListReq): Promise<BookmarkList>;
}

export const BookmarksServiceName = "meme.api.Bookmarks";
export class BookmarksClientImpl implements Bookmarks {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || BookmarksServiceName;
    this.rpc = rpc;
    this.Add = this.Add.bind(this);
    this.Remove = this.Remove.bind(this);
    this.List = this.List.bind(this);
  }
  Add(request: BookmarksAddReq): Promise<Void> {
    const data = BookmarksAddReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Add", data);
    return promise.then((data) => Void.fromJSON(data));
  }

  Remove(request: BookmarksAddReq): Promise<Void> {
    const data = BookmarksAddReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Remove", data);
    return promise.then((data) => Void.fromJSON(data));
  }

  List(request: BookmarkListReq): Promise<BookmarkList> {
    const data = BookmarkListReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "List", data);
    return promise.then((data) => BookmarkList.fromJSON(data));
  }
}

export interface Photos {
  Upload(request: UploadReq): Promise<UploadResp>;
}

export const PhotosServiceName = "meme.api.Photos";
export class PhotosClientImpl implements Photos {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || PhotosServiceName;
    this.rpc = rpc;
    this.Upload = this.Upload.bind(this);
  }
  Upload(request: UploadReq): Promise<UploadResp> {
    const data = UploadReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Upload", data);
    return promise.then((data) => UploadResp.fromJSON(data));
  }
}

export interface ImageProxy {
  Resize(request: ResizeReq): Promise<ResizeResp>;
}

export const ImageProxyServiceName = "meme.api.ImageProxy";
export class ImageProxyClientImpl implements ImageProxy {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || ImageProxyServiceName;
    this.rpc = rpc;
    this.Resize = this.Resize.bind(this);
  }
  Resize(request: ResizeReq): Promise<ResizeResp> {
    const data = ResizeReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "Resize", data);
    return promise.then((data) => ResizeResp.fromJSON(data));
  }
}

export interface Realtime {
  GetEvents(request: GetEventsReq): Promise<GetEventsResp>;
  SendEvent(request: SendEventReq): Promise<Void>;
}

export const RealtimeServiceName = "meme.api.Realtime";
export class RealtimeClientImpl implements Realtime {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || RealtimeServiceName;
    this.rpc = rpc;
    this.GetEvents = this.GetEvents.bind(this);
    this.SendEvent = this.SendEvent.bind(this);
  }
  GetEvents(request: GetEventsReq): Promise<GetEventsResp> {
    const data = GetEventsReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "GetEvents", data);
    return promise.then((data) => GetEventsResp.fromJSON(data));
  }

  SendEvent(request: SendEventReq): Promise<Void> {
    const data = SendEventReq.toJSON(request);
    //@ts-ignore-line
    const promise = this.rpc.request(this.service, "SendEvent", data);
    return promise.then((data) => Void.fromJSON(data));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(globalThis.String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
