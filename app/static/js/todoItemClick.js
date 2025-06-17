document.querySelectorAll('.todo-item').forEach(item => {
    item.addEventListener('click', () => {
    const id = item.dataset.id;  // data-id属性からIDを取得
    window.location.href = `/todo/${id}`;
    });
});
