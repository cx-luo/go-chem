/****************************************************************************
 * Copyright (C) from 2009 to Present EPAM Systems.
 *
 * This file is part of Indigo toolkit.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 ***************************************************************************/

package reaction

// Rect2f represents a 2D rectangle with floating-point coordinates
type Rect2f struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// SpecialCondition represents a special condition in a reaction
type SpecialCondition struct {
	MetaIdx int
	BBox    Rect2f
}

// NewSpecialCondition creates a new SpecialCondition
func NewSpecialCondition(idx int, box Rect2f) SpecialCondition {
	return SpecialCondition{
		MetaIdx: idx,
		BBox:    box,
	}
}

// AromaticityOptions represents options for aromatization
type AromaticityOptions struct {
	Method        string
	UniqueDoublet bool
}

// MetaDataStorage represents metadata storage
type MetaDataStorage struct {
	data map[string]interface{}
}

// NewMetaDataStorage creates a new MetaDataStorage
func NewMetaDataStorage() *MetaDataStorage {
	return &MetaDataStorage{
		data: make(map[string]interface{}),
	}
}

// Clone creates a copy of the metadata storage
func (m *MetaDataStorage) Clone(other *MetaDataStorage) {
	m.data = make(map[string]interface{})
	for k, v := range other.data {
		m.data[k] = v
	}
}

// GetMetaCount returns the count of metadata for a given ID
func (m *MetaDataStorage) GetMetaCount(id string) int {
	if val, ok := m.data[id]; ok {
		if arr, ok := val.([]interface{}); ok {
			return len(arr)
		}
	}
	return 0
}

// PropertiesMap represents a map of properties
type PropertiesMap struct {
	properties map[string]string
}

// NewPropertiesMap creates a new PropertiesMap
func NewPropertiesMap() *PropertiesMap {
	return &PropertiesMap{
		properties: make(map[string]string),
	}
}

// Copy copies properties from another map
func (p *PropertiesMap) Copy(other *PropertiesMap) {
	p.properties = make(map[string]string)
	for k, v := range other.properties {
		p.properties[k] = v
	}
}

// Get retrieves a property value
func (p *PropertiesMap) Get(key string) (string, bool) {
	val, ok := p.properties[key]
	return val, ok
}

// Set sets a property value
func (p *PropertiesMap) Set(key, value string) {
	p.properties[key] = value
}
