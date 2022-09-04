import React from "react";
import {BrowseResult, User, UserPage} from "../store/types";
import {Link} from "./Link";

export function PostUser(props: { user: User }) {
    const [isVisible, setIsVisible] = React.useState(false);
    const [userData, setUserData] = React.useState<BrowseResult>();
    const [isUserLoaded, setIsUserLoaded] = React.useState(false);

    let className = "UserNamePopup";
    if (!isVisible) {
        className += " UserNamePopup__hidden";
    }

    const onMouseEnter = () => {
        setIsVisible(true);

        if (isUserLoaded) return;

        setIsUserLoaded(true);
        fetch("http://localhost:8000/browse?url=/users/" + props.user.id)
            .then(r => r.json())
            .then(r => setUserData(r))
    }

    let userDetails = '...LOADING...';
    if (userData && userData.componentData) {
        const pageUser = userData.componentData as UserPage;
        const userObject = pageUser.nodes?.users?.find(item => item.id == pageUser.pageUser);
        if (userObject) {
            userDetails = "Name: " + userObject.name + ", posts: " + pageUser.posts?.length;
        }
    }

    return (
        <div className="UserName"
             onMouseEnter={onMouseEnter}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>
            From: <Link href={props.user.href}>{props.user.name}</Link>
        </div>
    )
}
