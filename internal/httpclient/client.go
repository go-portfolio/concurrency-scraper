package httpclient

import (
	"io/ioutil"
	"net/http"
	"time"
)

// Client — интерфейс для HTTP-клиента.
// Он определяет метод Fetch, который загружает содержимое страницы по URL.
type Client interface {
	Fetch(url string) (string, error)
}

// httpClient — конкретная реализация интерфейса Client.
// Внутри хранит стандартный http.Client с настройками.
type httpClient struct {
	client *http.Client
}

// New создаёт новый httpClient с таймаутом 10 секунд.
// Возвращает объект, который реализует интерфейс Client.
func New() Client {
	return &httpClient{
		client: &http.Client{
			Timeout: 10 * time.Second, // ограничиваем время выполнения запроса
		},
	}
}

// Fetch делает GET-запрос по указанному URL.
// Возвращает содержимое ответа как строку или ошибку, если запрос не удался.
func (c *httpClient) Fetch(url string) (string, error) {
	// Выполняем HTTP-запрос
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	// Гарантируем, что соединение будет закрыто после чтения
	defer resp.Body.Close()

	// Читаем тело ответа полностью
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Возвращаем результат в виде строки
	return string(data), nil
}
