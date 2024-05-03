package filler

import (
	model "github.com/askaroe/social-media-api/pkg/social-media/model"
	_ "github.com/lib/pq"
)

func PopulateDatabase(models model.Models) error {
	for _, post := range posts {
		models.Posts.Insert(&post)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

var posts = []model.Post{
	{Image: "image1.jpg", Caption: "First post", UserId: "user1"},
	{Image: "image2.jpg", Caption: "Second post", UserId: "user2"},
	{Image: "image3.jpg", Caption: "Third post", UserId: "user3"},
	{Image: "image4.jpg", Caption: "Fourth post", UserId: "user4"},
	{Image: "image5.jpg", Caption: "Fifth post", UserId: "user5"},
	{Image: "image6.jpg", Caption: "Sixth post", UserId: "user6"},
	{Image: "image7.jpg", Caption: "Seventh post", UserId: "user7"},
	{Image: "image8.jpg", Caption: "Eighth post", UserId: "user8"},
	{Image: "image9.jpg", Caption: "Ninth post", UserId: "user9"},
	{Image: "image10.jpg", Caption: "Tenth post", UserId: "user10"},
}

var users = []model.User{
	{ProfilePhoto: "profile1.jpg", Name: "John Doe", Username: "johndoe", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", Email: "john@example.com", Password: "password1", Age: "30"},
	{ProfilePhoto: "profile2.jpg", Name: "Jane Smith", Username: "janesmith", Description: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", Email: "jane@example.com", Password: "password2", Age: "25"},
	{ProfilePhoto: "profile3.jpg", Name: "Alice Johnson", Username: "alicejohnson", Description: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.", Email: "alice@example.com", Password: "password3", Age: "35"},
	{ProfilePhoto: "profile4.jpg", Name: "Bob Williams", Username: "bobwilliams", Description: "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.", Email: "bob@example.com", Password: "password4", Age: "28"},
	{ProfilePhoto: "profile5.jpg", Name: "Emily Brown", Username: "emilybrown", Description: "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", Email: "emily@example.com", Password: "password5", Age: "32"},
	{ProfilePhoto: "profile6.jpg", Name: "Michael Wilson", Username: "michaelwilson", Description: "Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur.", Email: "michael@example.com", Password: "password6", Age: "40"},
	{ProfilePhoto: "profile7.jpg", Name: "Emma Jones", Username: "emmajones", Description: "Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit.", Email: "emma@example.com", Password: "password7", Age: "27"},
	{ProfilePhoto: "profile8.jpg", Name: "Daniel Martinez", Username: "danielmartinez", Description: "Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat.", Email: "daniel@example.com", Password: "password8", Age: "33"},
	{ProfilePhoto: "profile9.jpg", Name: "Olivia Taylor", Username: "oliviataylor", Description: "Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet ut et voluptates repudiandae sint et molestiae non recusandae.", Email: "olivia@example.com", Password: "password9", Age: "29"},
	{ProfilePhoto: "profile10.jpg", Name: "William Anderson", Username: "williamanderson", Description: "Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat.", Email: "william@example.com", Password: "password10", Age: "31"},
}
