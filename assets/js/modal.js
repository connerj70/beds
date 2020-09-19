function modal(content, closeModal, submitModal, addingFriend, addFriend) {

    var e = React.createElement

    let showAddFriend
    if (addingFriend !== undefined) {
        showAddFriend = e('div', null,
            addingFriend.email,
            e('button', { onClick: addFriend }, 'Add Friend')
        )
    }

    return React.createElement(
        'div',
        { className: "c-modal" },
        content,
        React.createElement(
            'button',
            { onClick: closeModal },
            'close'
        ),
        React.createElement(
            'button',
            { onClick: submitModal },
            'submit'
        ),
        e(
            'div',
            null,
            showAddFriend
        )
    )
}

export default modal