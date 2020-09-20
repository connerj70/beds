'use strict';

import modal from "./modal.js";
import getUserIDFromCookies from "./cookies.js"

const e = React.createElement;

class FriendsContainer extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            friends: [],
            showAddFriend: false,
            searchTerm: "",
            addingFriend: {}
        };

        this.closeModal = this.closeModal.bind(this);
        this.submitModal = this.submitModal.bind(this);
        this.handleSearchChange = this.handleSearchChange.bind(this);
        this.addFriend = this.addFriend.bind(this);
    }

    componentDidMount() {
        let userID = getUserIDFromCookies(document.cookie)
        fetch(`/friends/list/${userID}`).then(resp => {
            resp.json().then(data => {
                this.setState({
                    friends: data
                })
            })
        })
    }

    closeModal() {
        this.setState({
            showAddFriend: false
        })
    }

    submitModal() {
        // go query the database for a user with this email
        fetch("/users/by_email", { method: "POST", headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email: this.state.searchTerm }) }).then(resp => {
            if (resp.status !== 200) {
                resp.json().then(err => {
                    alert(`There was an error with your request ${err.error}`)
                })
                return
            }
            resp.json().then(data => {
                this.setState({
                    addingFriend: data
                })
            })
        })
    }

    handleSearchChange(val) {
        this.setState({
            searchTerm: val
        })
    }

    addFriend() {
        let userID = getUserIDFromCookies(document.cookie)

        // make a friend request to the database.
        fetch("/friends/create", { method: "POST", headers: { 'Content-Type': "application/json" }, body: JSON.stringify({ requester_id: userID, receiver_id: this.state.addingFriend.id }) })
            .then(resp => {
                if (resp.status !== 200) {
                    resp.json().then(err => {
                        alert(`there was an error in your request: ${err.error}`)
                        return
                    })
                }
                resp.json().then(data => {
                    console.log(data)
                    this.setState({
                        friends: this.state.friends.concat(data)
                    })
                })
            })
    }

    render() {
        let approvedFriends = this.state.friends.map(friend => {
            if (friend.approved) {
                return e(
                    'div',
                    { "key": friend.email },
                    friend.email
                )
            }
        });

        let pendingFriends = this.state.friends.map(friend => {
            if (!friend.approved) {
                return e(
                    'div',
                    { "key": friend.email },
                    friend.email
                )
            }
        });

        let addFriend;
        if (this.state.showAddFriend) {
            addFriend = modal(e('div', null, e('input', { type: "text", placeholder: "search for friends...", value: this.state.searchTerm, onChange: (e) => this.handleSearchChange(e.target.value) })), this.closeModal, this.submitModal, this.state.addingFriend, this.addFriend)
        }

        return e(
            'div',
            null,
            e('div', null, addFriend),
            e('button', { onClick: () => { this.setState({ showAddFriend: !this.state.showAddFriend }) } }, "Add Friend"),
            e('div', null, 'Approved Friends', approvedFriends),
            e('div', null, 'Pending Friends', pendingFriends)
        );
    }
}

const domContainer = document.querySelector('#friends-container');
ReactDOM.render(e(FriendsContainer), domContainer);