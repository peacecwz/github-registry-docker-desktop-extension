const updateStore = (state, payload) => {
    return {
        ...state,
        github: {
            ...state.github,
            ...payload,
        },
    };
}

export default updateStore;