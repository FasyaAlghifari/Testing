import axios from "axios";

const API_URL = "http://localhost:8080/notifications";

export const GetNotification = (callback) => {
  return axios
    .get(`${API_URL}`)
    .then((response) => {
      callback(response.data); // ensure this matches the structure sent from the backend
    })
    .catch((error) => {
      throw new Error(`Gagal menampilkan notif. Alasan: ${error.message}`);
    });
};

export function deleteNotification(id) {
  return axios
    .delete(`${API_URL}/${id}`)
    .then((response) => {
      return response.data;
    })
    .catch((error) => {
      throw new Error(`Gagal menghapus notif. Alasan: ${error.message}`);
    });
}

// Tambahkan fungsi untuk menghubungkan WebSocket
export const connectWebSocket = (callback) => {
  const ws = new WebSocket('ws://localhost:8080/ws'); // Sesuaikan URL dengan endpoint WebSocket Anda

  ws.onmessage = (event) => {
    const notification = JSON.parse(event.data);
    callback(notification); // Panggil callback dengan data notifikasi yang diterima
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  ws.onclose = () => {
    console.log('WebSocket connection closed');
  };

  return ws; // Mengembalikan WebSocket connection untuk kemungkinan penggunaan lebih lanjut
};
