package db

func (d *DB) GetTownByID(townID string) (Town, error) {
	query := `SELECT id, name, feed_url FROM towns WHERE id = ?;`

	rows, err := d.db.Query(query, townID)
	if err != nil {
		return Town{}, err
	}
	defer rows.Close()

	var town Town

	for rows.Next() {
		err := rows.Scan(&town.ID,
			&town.Name, &town.FeedUrl)
		if err != nil {
			return Town{}, err
		}
	}

	return town, nil
}

func (d *DB) ListTowns() ([]Town, error) {
	query := `SELECT id, name, feed_url FROM towns ORDER BY name;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var towns []Town

	for rows.Next() {
		var town Town
		err := rows.Scan(&town.ID,
			&town.Name, &town.FeedUrl)
		if err != nil {
			return nil, err
		}

		towns = append(towns, town)
	}

	return towns, nil
}

func (d *DB) AddSubscription(chatId int64, townID string) error {
	query := `
		INSERT INTO subscriptions (chat_id, town) VALUES (?, ?);
	`

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatId, townID)
	return err
}

func (d *DB) GetSubscriptionsByChatId(chatId int64) ([]Town, error) {
	query := `
	select t.id, t.name, t.feed_url FROM subscriptions
	JOIN towns AS t
	ON town = t.id
	WHERE chat_id = ?;
`

	rows, err := d.db.Query(query, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var towns []Town

	for rows.Next() {
		var town Town
		err := rows.Scan(&town.ID,
			&town.Name, &town.FeedUrl)
		if err != nil {
			return nil, err
		}

		towns = append(towns, town)
	}

	return towns, nil
}

func (d *DB) RemoveSubscription(chatId int64, townID string) error {
	query := `
	DELETE FROM subscriptions WHERE chat_id = ? AND town = ?;
`

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatId, townID)
	return err
}

func (d *DB) GetSubscribedTowns() ([]Town, error) {
	query := `
SELECT DISTINCT id, name, feed_url, last_updated FROM towns JOIN subscriptions AS s ON s.town = id;
`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var towns []Town

	for rows.Next() {
		var town Town
		err := rows.Scan(&town.ID, &town.Name, &town.FeedUrl, &town.LastUpdated)
		if err != nil {
			return nil, err
		}

		towns = append(towns, town)
	}

	return towns, nil
}

func (d *DB) UpdateLastUpdated(towndId string, lastUpdated int64) error {
	query := `
	UPDATE towns SET last_updated = ?  WHERE id = ?;
	`

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(lastUpdated, towndId)

	return nil
}

func (d *DB) GetSubscriberChatId(townId string) ([]int64, error) {
	query := `
SELECT DISTINCT chat_id FROM subscriptions WHERE town = ?;
`

	rows, err := d.db.Query(query, townId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIds []int64

	for rows.Next() {
		var chatId int64
		err := rows.Scan(&chatId)
		if err != nil {
			return nil, err
		}

		chatIds = append(chatIds, chatId)
	}

	return chatIds, nil
}
