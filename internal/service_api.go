package internal

import (
	"context"
	"fmt"

	"github.com/MustafaLo/noted/models"
	"github.com/dstotijn/go-notion"
)

func FilterDatabase(s *models.APIService, DB_ID string, query string)(*notion.DatabaseQueryResponse, error){
	fmt.Println("Retrieving Notes...")
	database_query := createDatabaseQuery(query)
	query_response, err := s.Client.QueryDatabase(context.Background(), DB_ID, &database_query)
	if err != nil{
		return nil, fmt.Errorf("error retrieving results %w", err)
	}
	return &query_response, nil
}

func CreateDatabaseEntry(s *models.APIService, DB_ID string, fileMetaData models.FileMetadata, note string, lines string, category string)(string, error){
	fmt.Println("Creating Page...")
	database_page_properties := createDatabasePageProperties(fileMetaData, note, lines, category)
	page_parameters := createPageParams(DB_ID, database_page_properties)
	page, err := s.Client.CreatePage(context.Background(), page_parameters)
	if err != nil{
		return "", fmt.Errorf("error creating page: %w", err)
	}

	fmt.Println("\nPage created successfully: ", page.URL)
	return page.ID, nil
}

func UpdateDatabaseEntry(s *models.APIService, PAGE_ID string, code string, language string, note string)(error){
	code_heading_block := createHeadingBlock("Code Snippet")
	code_block := createCodeBlock(code, language)
	note_block := createNoteBlock(note)
	note_heading_block := createHeadingBlock("Note")
	_, err := s.Client.AppendBlockChildren(context.Background(), PAGE_ID, []notion.Block{code_heading_block, code_block, note_heading_block, note_block})
	if err != nil {
		return fmt.Errorf("failed to append code block: %w", err)
	}
	return nil
}

func createDatabaseQuery(filter string)(notion.DatabaseQuery){
	return notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "File Name",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Title: &notion.TextPropertyFilter{
					Equals: filter,
				},
			},
		},
	}
}

func createHeadingBlock(content string)(notion.Heading2Block){
	return notion.Heading2Block{
		RichText: []notion.RichText{{Text: &notion.Text{Content: content}}},
	}
}

func createNoteBlock(content string)(notion.ParagraphBlock){
	return notion.ParagraphBlock{
		RichText: []notion.RichText{{Text: &notion.Text{Content: content}}},
	}
}

func createCodeBlock(content string, language string)(notion.CodeBlock){
	return notion.CodeBlock{
		RichText: []notion.RichText{{Text: &notion.Text{Content: content}}},
		Language: &language,
	}
}


func createPageParams(DB_ID string, db_page_props notion.DatabasePageProperties)(notion.CreatePageParams){
	return notion.CreatePageParams{
		ParentType: "database_id",
		ParentID: DB_ID,
		DatabasePageProperties: &db_page_props,
	}
}

func createDatabasePageProperties(fileMetaData models.FileMetadata, note string, lines string, category string)(notion.DatabasePageProperties){
	return notion.DatabasePageProperties{
		"File Name": notion.DatabasePageProperty{
			Title: []notion.RichText{{Text:&notion.Text {Content: fileMetaData.FileName}}},
		},
		"Note": notion.DatabasePageProperty{
			RichText: []notion.RichText{{Text: &notion.Text{Content: note}}},
		},
		"Line Numbers": notion.DatabasePageProperty{
			RichText: []notion.RichText{{Text: &notion.Text{Content: lines}}},
		},
		"Category": notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: category},
		},
	}
}
