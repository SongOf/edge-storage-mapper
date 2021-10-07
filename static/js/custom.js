function deleteNote(noteId) {
    let choice = confirm("Are you sure to delete this note ?");

    if (choice == true) {
        location.href = "notepad/delete/" + noteId;
        return false
    }
}

function deleteDevice(ip) {
    let choice = confirm("Are you sure to delete this device ?");
    console.log(ip)
    if (choice == true) {
        location.href = "device/delete/" + ip
        return false
    }
}